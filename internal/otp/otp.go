package otp

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/otp/config"
	"github.com/hoitek/Maja-Service/internal/otp/handlers"
	"github.com/hoitek/Maja-Service/internal/otp/ports"
	"github.com/hoitek/Maja-Service/internal/otp/repositories"
	"github.com/hoitek/Maja-Service/internal/otp/service"
	userPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the otp domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the otp domain module
var Module = &module{}

// GetService returns a new instance of the otp service
func (m *module) GetService(pDB *sql.DB) ports.OTPService {
	// otp repository database based on the environment
	otpRepositoryPostgresDB := exp.TerIf[ports.OTPRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewOTPRepositoryStub(),
		repositories.NewOTPRepositoryPostgresDB(pDB),
	)
	// otp service and inject the otp repository database and grpc
	otpService := service.NewOTPService(otpRepositoryPostgresDB)
	return otpService
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
	config.OTPConfig = &c
	return m
}

// SetDatabases sets the database for the staff domain module
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

// RegisterHTTP registers the otp domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.OTPHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.OTPHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for otp module
	r.Use(json.Middleware)

	// Create a new otp handler
	handler, err := handlers.NewOTPHandler(r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.OTPHandler{}, err
	}

	// Return the otp handler
	return handler, nil
}
