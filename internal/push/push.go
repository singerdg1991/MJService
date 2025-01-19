package push

import (
	"database/sql"
	"errors"

	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/push/config"
	"github.com/hoitek/Maja-Service/internal/push/handlers"
	"github.com/hoitek/Maja-Service/internal/push/ports"
	"github.com/hoitek/Maja-Service/internal/push/repositories"
	"github.com/hoitek/Maja-Service/internal/push/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the push domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the push domain module
var Module = &module{}

// GetService returns a new instance of the push service
func (m *module) GetService(pDB *sql.DB) ports.PushService {
	// push repository database based on the environment
	pushRepositoryPostgresDB := exp.TerIf[ports.PushRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewPushRepositoryStub(),
		repositories.NewPushRepositoryPostgresDB(pDB),
	)
	// push service and inject the push repository database and grpc
	pushService := service.NewPushService(pushRepositoryPostgresDB, m.MinIOStorage)
	return &pushService
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
	config.PushConfig = &c
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

// RegisterHTTP registers the push domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.PushHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.PushHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for push module
	r.Use(json.Middleware)

	// Create a new push handler
	handler, err := handlers.NewPushHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.PushHandler{}, err
	}

	// Return the push handler
	return handler, nil
}
