package todo

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
	"github.com/hoitek/Maja-Service/internal/todo/config"
	"github.com/hoitek/Maja-Service/internal/todo/handlers"
	"github.com/hoitek/Maja-Service/internal/todo/ports"
	"github.com/hoitek/Maja-Service/internal/todo/repositories"
	"github.com/hoitek/Maja-Service/internal/todo/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the todos domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the todos domain module
var Module = &module{}

// GetService returns a new instance of the todos service
func (m *module) GetService(pDB *sql.DB) ports.TodoService {
	// todos repository database based on the environment
	todoRepositoryPostgresDB := exp.TerIf[ports.TodoRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewTodoRepositoryStub(),
		repositories.NewTodoRepositoryPostgresDB(pDB),
	)
	// todos service and inject the todos repository database and grpc
	todoService := service.NewTodoService(todoRepositoryPostgresDB, m.MinIOStorage)
	return &todoService
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
	config.TodoConfig = &c
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

// RegisterHTTP registers the todos domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.TodoHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.TodoHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for todos module
	r.Use(json.Middleware)

	// Create a new todos handler
	handler, err := handlers.NewTodoHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.TodoHandler{}, err
	}

	// Return the todos handler
	return handler, nil
}
