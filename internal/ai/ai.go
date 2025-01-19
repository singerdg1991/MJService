package ai

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/ai/config"
	"github.com/hoitek/Maja-Service/internal/ai/handlers"
	"github.com/hoitek/Maja-Service/internal/ai/ports"
	service "github.com/hoitek/Maja-Service/internal/ai/services"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the ai domain
type module struct {
	Config       config.ConfigType
	MinIOStorage *storage.MinIO
	MongoDB      *mongo.Client
	PostgresDB   *sql.DB
}

// Module is a global variable for the ai domain module
var Module = &module{}

// GetService returns a new instance of the ai service
func (m *module) GetService() ports.AIService {
	// ai service and inject the ai repository database and grpc
	aiService := service.NewAIService(m.MinIOStorage)
	return &aiService
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
	config.AIConfig = &c
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// SetDatabases sets the databases
func (m *module) SetDatabases(pDB *sql.DB, mDB *mongo.Client) *module {
	m.PostgresDB = pDB
	m.MongoDB = mDB
	return m
}

// RegisterHTTP registers the ai domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.AIHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.AIHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for ai module
	r.Use(json.Middleware)

	// Create a new ai handler
	handler, err := handlers.NewAIHandler(
		r,
		m.GetService(),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.AIHandler{}, err
	}

	// Return the ai handler
	return handler, nil
}
