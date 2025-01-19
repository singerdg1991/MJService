package warning

import (
	"database/sql"
	"errors"

	rPorts "github.com/hoitek/Maja-Service/internal/punishment/ports"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hoitek/Maja-Service/storage"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	rRepositories "github.com/hoitek/Maja-Service/internal/punishment/repositories"
	rService "github.com/hoitek/Maja-Service/internal/punishment/service"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/config"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/handlers"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/ports"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/repositories"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the warning domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the warning domain module
var Module = &module{}

// GetService returns a new instance of the warning service
func (m *module) GetService(pDB *sql.DB) ports.WarningService {
	// warning repository database based on the environment
	warningRepositoryPostgresDB := exp.TerIf[ports.WarningRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewWarningRepositoryStub(),
		repositories.NewWarningRepositoryPostgresDB(pDB),
	)
	// warning service and inject the warning repository database and grpc
	warningService := service.NewWarningService(warningRepositoryPostgresDB, m.MinIOStorage)
	return &warningService
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

// GetPunishmentService returns a new instance of the punishment service
func (m *module) GetPunishmentService(pDB *sql.DB) rPorts.PunishmentService {
	// punishment repository database based on the environment
	punishmentRepositoryPostgresDB := exp.TerIf[rPorts.PunishmentRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		rRepositories.NewPunishmentRepositoryStub(),
		rRepositories.NewPunishmentRepositoryPostgresDB(pDB),
	)
	// punishment service and inject the punishment repository database and grpc
	punishmentService := rService.NewPunishmentService(punishmentRepositoryPostgresDB, m.MinIOStorage)
	return &punishmentService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.WarningConfig = &c
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

// RegisterHTTP registers the warning domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.WarningHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.WarningHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for warning module
	r.Use(json.Middleware)

	// Create a new warning handler
	handler, err := handlers.NewWarningHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetPunishmentService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.WarningHandler{}, err
	}

	// Return the warning handler
	return handler, nil
}
