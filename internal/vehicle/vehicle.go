package vehicle

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	cPorts "github.com/hoitek/Maja-Service/internal/company/ports"
	cRepositories "github.com/hoitek/Maja-Service/internal/company/repositories"
	cService "github.com/hoitek/Maja-Service/internal/company/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/internal/vehicle/config"
	"github.com/hoitek/Maja-Service/internal/vehicle/handlers"
	"github.com/hoitek/Maja-Service/internal/vehicle/ports"
	"github.com/hoitek/Maja-Service/internal/vehicle/repositories"
	"github.com/hoitek/Maja-Service/internal/vehicle/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the vehicle domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the vehicle domain module
var Module = &module{}

// GetService returns a new instance of the vehicle service
func (m *module) GetService(pDB *sql.DB) ports.VehicleService {
	// vehicle repository database based on the environment
	vehicleRepositoryPostgresDB := exp.TerIf[ports.VehicleRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewVehicleRepositoryStub(),
		repositories.NewVehicleRepositoryPostgresDB(pDB),
	)
	// vehicle service and inject the vehicle repository database and grpc
	vehicleService := service.NewVehicleService(vehicleRepositoryPostgresDB, m.MinIOStorage)
	return &vehicleService
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

// GetCompanyService returns a new instance of the company service
func (m *module) GetCompanyService(pDB *sql.DB) cPorts.CompanyService {
	// company repository database based on the environment
	companyRepositoryPostgresDB := exp.TerIf[cPorts.CompanyRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		cRepositories.NewCompanyRepositoryStub(),
		cRepositories.NewCompanyRepositoryPostgresDB(pDB),
	)
	// company service and inject the company repository database and grpc
	companyService := cService.NewCompanyService(companyRepositoryPostgresDB, m.MinIOStorage)
	return &companyService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.VehicleConfig = &c
	return m
}

// SetDatabases sets the database for the user domain module
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

// RegisterHTTP registers the vehicle domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.VehicleHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.VehicleHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for vehicle module
	r.Use(json.Middleware)

	// Create a new vehicle handler
	handler, err := handlers.NewVehicleHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
		m.GetCompanyService(m.PostgresDB),
	)
	if err != nil {
		return handlers.VehicleHandler{}, err
	}

	// Return the vehicle handler
	return handler, nil
}
