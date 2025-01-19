package ticket

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
	"github.com/hoitek/Maja-Service/internal/ticket/config"
	"github.com/hoitek/Maja-Service/internal/ticket/handlers"
	"github.com/hoitek/Maja-Service/internal/ticket/ports"
	"github.com/hoitek/Maja-Service/internal/ticket/repositories"
	"github.com/hoitek/Maja-Service/internal/ticket/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the ticketCategory domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the ticketCategory domain module
var Module = &module{}

// GetService returns a new instance of the ticketCategory service
func (m *module) GetService(pDB *sql.DB) ports.TicketService {
	// ticketCategory repository database based on the environment
	ticketCategoryRepositoryPostgresDB := exp.TerIf[ports.TicketRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewTicketRepositoryStub(),
		repositories.NewTicketRepositoryPostgresDB(pDB),
	)
	// ticketCategory service and inject the ticketCategory repository database and grpc
	ticketCategoryService := service.NewTicketService(ticketCategoryRepositoryPostgresDB, m.MinIOStorage)
	return &ticketCategoryService
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
	config.TicketConfig = &c
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

// RegisterHTTP registers the ticketCategory domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.TicketHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.TicketHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for ticketCategory module
	r.Use(json.Middleware)

	// Create a new ticketCategory handler
	handler, err := handlers.NewTicketHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.TicketHandler{}, err
	}

	// Return the ticketCategory handler
	return handler, nil
}
