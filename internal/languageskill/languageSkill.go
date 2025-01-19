package languageskill

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/languageskill/config"
	"github.com/hoitek/Maja-Service/internal/languageskill/handlers"
	"github.com/hoitek/Maja-Service/internal/languageskill/ports"
	"github.com/hoitek/Maja-Service/internal/languageskill/repositories"
	"github.com/hoitek/Maja-Service/internal/languageskill/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the languageskill domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the languageskill domain module
var Module = &module{}

// GetService returns a new instance of the languageskill service
func (m *module) GetService(pDB *sql.DB) ports.LanguageSkillService {
	// languageskill repository database based on the environment
	languageskillRepositoryPostgresDB := exp.TerIf[ports.LanguageSkillRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewLanguageSkillRepositoryStub(),
		repositories.NewLanguageSkillRepositoryPostgresDB(pDB),
	)
	// languageskill service and inject the languageskill repository database and grpc
	languageskillService := service.NewLanguageSkillService(languageskillRepositoryPostgresDB, m.MinIOStorage)
	return &languageskillService
}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB, mDB *mongo.Client) uPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[uPorts.UserRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		uRepositories.NewUserRepositoryStub(),
		uRepositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := uRepositories.NewUserRepositoryMongoDB(mDB)

	// user service and inject the user repository database and grpc
	userService := uService.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, storage.MinIOStorage)
	return userService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.LanguageSkillConfig = &c
	return m
}

// SetDatabases sets the databases
func (m *module) SetDatabases(pDB *sql.DB, mDB *mongo.Client) *module {
	m.PostgresDB = pDB
	m.MongoDB = mDB
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// RegisterHTTP registers the languageskill domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.LanguageSkillHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.LanguageSkillHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for languageskill module
	r.Use(json.Middleware)

	// Create a new languageskill handler
	handler, err := handlers.NewLanguageSkillHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.LanguageSkillHandler{}, err
	}

	// Return the languageskill handler
	return handler, nil
}
