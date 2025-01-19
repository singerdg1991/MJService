package notification

import (
	"database/sql"
	"errors"
	userPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hoitek/Maja-Service/storage"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/notification/config"
	"github.com/hoitek/Maja-Service/internal/notification/handlers"
	"github.com/hoitek/Maja-Service/internal/notification/ports"
	"github.com/hoitek/Maja-Service/internal/notification/repositories"
	"github.com/hoitek/Maja-Service/internal/notification/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the notification domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the notification domain module
var Module = &module{}

// GetService returns a new instance of the notification service
func (m *module) GetService(pDB *sql.DB) ports.NotificationService {
	// notification repository database based on the environment
	notificationRepositoryPostgresDB := exp.TerIf[ports.NotificationRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewNotificationRepositoryStub(),
		repositories.NewNotificationRepositoryPostgresDB(pDB),
	)
	// notification service and inject the notification repository database and grpc
	notificationService := service.NewNotificationService(notificationRepositoryPostgresDB, m.MinIOStorage)
	return &notificationService
}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB) userPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[userPorts.UserRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		uRepositories.NewUserRepositoryStub(),
		uRepositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := uRepositories.NewUserRepositoryMongoDB(m.MongoDB)

	// user service and inject the user repository database and grpc
	userService := uService.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, m.MinIOStorage)
	return userService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.NotificationConfig = &c
	return m
}

// SetDatabases sets the database for the todos domain module
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

// RegisterHTTP registers the notification domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.NotificationHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.NotificationHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for notification module
	r.Use(json.Middleware)

	// Create a new notification handler
	handler, err := handlers.NewNotificationHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.NotificationHandler{}, err
	}

	// Return the notification handler
	return handler, nil
}
