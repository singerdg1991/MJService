package contracttype

import (
	"database/sql"
	"errors"

	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/contracttype/config"
	"github.com/hoitek/Maja-Service/internal/contracttype/handlers"
	"github.com/hoitek/Maja-Service/internal/contracttype/ports"
	"github.com/hoitek/Maja-Service/internal/contracttype/repositories"
	"github.com/hoitek/Maja-Service/internal/contracttype/service"
	"github.com/hoitek/Middlewares/json"

	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
)

// module is a module for the contractType domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the contractType domain module
var Module = &module{}

// GetService returns a new instance of the contractType service
func (m *module) GetService(pDB *sql.DB) ports.ContractTypeService {
	// contractType repository database based on the environment
	contractTypeRepositoryPostgresDB := exp.TerIf[ports.ContractTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewContractTypeRepositoryStub(),
		repositories.NewContractTypeRepositoryPostgresDB(pDB),
	)
	// ContractType service and inject the contractType repository database and grpc
	contractTypeService := service.NewContractTypeService(contractTypeRepositoryPostgresDB, m.MinIOStorage)
	return &contractTypeService
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
	config.ContractTypeConfig = &c
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

// RegisterHTTP registers the contractType domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ContractTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ContractTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for contractType module
	r.Use(json.Middleware)

	// Create a new contractType handler
	handler, err := handlers.NewContractTypeHandler(r, m.GetService(m.PostgresDB), m.GetUserService(m.PostgresDB, m.MongoDB))
	if err != nil {
		return handlers.ContractTypeHandler{}, err
	}

	// Return the contractType handler
	return handler, nil
}
