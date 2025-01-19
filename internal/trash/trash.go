package trash

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/trash/config"
	"github.com/hoitek/Maja-Service/internal/trash/handlers"
	"github.com/hoitek/Maja-Service/internal/trash/ports"
	"github.com/hoitek/Maja-Service/internal/trash/repositories"
	"github.com/hoitek/Maja-Service/internal/trash/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the trash domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the trash domain module
var Module = &module{}

// GetService returns a new instance of the trash service
func (m *module) GetService(pDB *sql.DB) ports.TrashService {
	// trash repository database based on the environment
	trashRepositoryPostgresDB := exp.TerIf[ports.TrashRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewTrashRepositoryStub(),
		repositories.NewTrashRepositoryPostgresDB(pDB),
	)
	// trash service and inject the trash repository database and grpc
	trashService := service.NewTrashService(trashRepositoryPostgresDB, m.MinIOStorage)
	return &trashService
}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB) uPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[uPorts.UserRepositoryPostgresDB](
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
	config.TrashConfig = &c
	return m
}

// SetDatabases sets the database for the trash domain module
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

// RegisterHTTP registers the trash domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.TrashHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.TrashHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for trash module
	r.Use(json.Middleware)

	// Create a new trash handler
	handler, err := handlers.NewTrashHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.TrashHandler{}, err
	}

	// Return the trash handler
	return handler, nil
}
