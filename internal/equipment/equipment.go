package equipment

import (
	"database/sql"
	"errors"

	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/equipment/config"
	"github.com/hoitek/Maja-Service/internal/equipment/handlers"
	"github.com/hoitek/Maja-Service/internal/equipment/ports"
	"github.com/hoitek/Maja-Service/internal/equipment/repositories"
	"github.com/hoitek/Maja-Service/internal/equipment/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the equipment domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the equipment domain module
var Module = &module{}

// GetService returns a new instance of the equipment service
func (m *module) GetService(pDB *sql.DB) ports.EquipmentService {
	// equipment repository database based on the environment
	equipmentRepositoryPostgresDB := exp.TerIf[ports.EquipmentRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewEquipmentRepositoryStub(),
		repositories.NewEquipmentRepositoryPostgresDB(pDB),
	)
	// equipment service and inject the equipment repository database and grpc
	equipmentService := service.NewEquipmentService(equipmentRepositoryPostgresDB, m.MinIOStorage)
	return &equipmentService
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
	config.EquipmentConfig = &c
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

// RegisterHTTP registers the equipment domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.EquipmentHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.EquipmentHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for equipment module
	r.Use(json.Middleware)

	// Create a new equipment handler
	handler, err := handlers.NewEquipmentHandler(r, m.GetService(m.PostgresDB), m.GetUserService(m.PostgresDB, m.MongoDB))
	if err != nil {
		return handlers.EquipmentHandler{}, err
	}

	// Return the equipment handler
	return handler, nil
}
