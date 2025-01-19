package oauth2

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/oauth2/config"
	"github.com/hoitek/Maja-Service/internal/oauth2/handlers"
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

// module is a module for the oauth2 domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the oauth2 domain module
var Module = &module{}

// GetService returns a new instance of the oauth2 service
func (m *module) GetService(pDB *sql.DB) ports.OTPService {
	// oauth2 repository database based on the environment
	oauth2RepositoryPostgresDB := exp.TerIf[ports.OTPRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewOTPRepositoryStub(),
		repositories.NewOTPRepositoryPostgresDB(pDB),
	)
	// oauth2 service and inject the oauth2 repository database and grpc
	oauth2Service := service.NewOTPService(oauth2RepositoryPostgresDB)
	return oauth2Service
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
	config.OAuth2Config = &c
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

// RegisterHTTP registers the oauth2 domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.OAuth2Handler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.OAuth2Handler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for oauth2 module
	r.Use(json.Middleware)

	// Create a new oauth2 handler
	handler, err := handlers.NewOAuth2Handler(r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.OAuth2Handler{}, err
	}

	// Return the oauth2 handler
	return handler, nil
}
