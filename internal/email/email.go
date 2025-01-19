package email

import (
	"database/sql"
	"errors"

	userPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/email/config"
	"github.com/hoitek/Maja-Service/internal/email/handlers"
	"github.com/hoitek/Maja-Service/internal/email/ports"
	"github.com/hoitek/Maja-Service/internal/email/repositories"
	"github.com/hoitek/Maja-Service/internal/email/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the email domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the email domain module
var Module = &module{}

// GetService returns a new instance of the email service
func (m *module) GetService(pDB *sql.DB) ports.EmailService {
	// email repository database based on the environment
	emailRepositoryPostgresDB := exp.TerIf[ports.EmailRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewEmailRepositoryStub(),
		repositories.NewEmailRepositoryPostgresDB(pDB),
	)
	// email service and inject the email repository database and grpc
	emailService := service.NewEmailService(emailRepositoryPostgresDB, m.MinIOStorage)
	return &emailService
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
	config.EmailConfig = &c
	return m
}

// SetDatabases sets the database for the email domain module
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

// RegisterHTTP registers the email domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.EmailHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.EmailHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for email module
	r.Use(json.Middleware)

	// Create a new email handler
	handler, err := handlers.NewEmailHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.EmailHandler{}, err
	}

	// Return the email handler
	return handler, nil
}
