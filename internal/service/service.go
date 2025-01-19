package service

import (
	"database/sql"
	"errors"

	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/service/config"
	"github.com/hoitek/Maja-Service/internal/service/handlers"
	"github.com/hoitek/Maja-Service/internal/service/ports"
	"github.com/hoitek/Maja-Service/internal/service/repositories"
	"github.com/hoitek/Maja-Service/internal/service/service"
	"github.com/hoitek/Middlewares/json"

	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
)

// module is a module for the service domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the service domain module
var Module = &module{}

// GetService returns a new instance of the service service
func (m *module) GetService(pDB *sql.DB) ports.ServiceService {
	// service repository database based on the environment
	serviceRepositoryPostgresDB := exp.TerIf[ports.ServiceRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewServiceRepositoryStub(),
		repositories.NewServiceRepositoryPostgresDB(pDB),
	)
	// service service and inject the service repository database and grpc
	serviceService := service.NewServiceService(serviceRepositoryPostgresDB, m.MinIOStorage)
	return &serviceService
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
	config.ServiceConfig = &c
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
	// Set the minio storage
	m.MinIOStorage = s
	return m
}

// RegisterHTTP registers the service domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ServiceHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ServiceHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for service module
	r.Use(json.Middleware)

	// Create a new service handler
	handler, err := handlers.NewServiceHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.ServiceHandler{}, err
	}

	// Return the service handler
	return handler, nil
}
