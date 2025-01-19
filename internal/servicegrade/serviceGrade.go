package servicegrade

import (
	"database/sql"
	"errors"

	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/servicegrade/config"
	"github.com/hoitek/Maja-Service/internal/servicegrade/handlers"
	"github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	"github.com/hoitek/Maja-Service/internal/servicegrade/repositories"
	"github.com/hoitek/Maja-Service/internal/servicegrade/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the servicegrade domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the servicegrade domain module
var Module = &module{}

// GetService returns a new instance of the servicegrade service
func (m *module) GetService(pDB *sql.DB) ports.ServiceGradeService {
	// servicegrade repository database based on the environment
	servicegradeRepositoryPostgresDB := exp.TerIf[ports.ServiceGradeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewServiceGradeRepositoryStub(),
		repositories.NewServiceGradeRepositoryPostgresDB(pDB),
	)
	// servicegrade service and inject the servicegrade repository database and grpc
	servicegradeService := service.NewServiceGradeService(servicegradeRepositoryPostgresDB, m.MinIOStorage)
	return &servicegradeService
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
	config.ServiceGradeConfig = &c
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

// RegisterHTTP registers the servicegrade domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ServiceGradeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ServiceGradeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for servicegrade module
	r.Use(json.Middleware)

	// Create a new servicegrade handler
	handler, err := handlers.NewServiceGradeHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.ServiceGradeHandler{}, err
	}

	// Return the servicegrade handler
	return handler, nil
}
