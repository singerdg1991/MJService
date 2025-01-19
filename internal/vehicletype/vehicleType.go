package vehicletype

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/vehicletype/config"
	"github.com/hoitek/Maja-Service/internal/vehicletype/handlers"
	"github.com/hoitek/Maja-Service/internal/vehicletype/ports"
	"github.com/hoitek/Maja-Service/internal/vehicletype/repositories"
	"github.com/hoitek/Maja-Service/internal/vehicletype/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the vehicletype domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the vehicletype domain module
var Module = &module{}

// GetService returns a new instance of the vehicletype service
func (m *module) GetService(pDB *sql.DB) service.VehicleTypeService {
	// vehicletype repository database based on the environment
	vehicletypeRepositoryPostgresDB := exp.TerIf[ports.VehicleTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewVehicleTypeRepositoryStub(),
		repositories.NewVehicleTypeRepositoryPostgresDB(pDB),
	)
	// vehicletype service and inject the vehicletype repository database and grpc
	vehicletypeService := service.NewVehicleTypeService(vehicletypeRepositoryPostgresDB, m.MinIOStorage)
	return vehicletypeService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.VehicleTypeConfig = &c
	return m
}

// SetDatabase sets the database for the setting domain module
func (m *module) SetDatabase(pDB *sql.DB) *module {
	m.PostgresDB = pDB
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// RegisterHTTP registers the vehicletype domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.VehicleTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.VehicleTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for vehicletype module
	r.Use(json.Middleware)

	// Create a new vehicletype handler
	handler, err := handlers.NewVehicleTypeHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.VehicleTypeHandler{}, err
	}

	// Return the vehicletype handler
	return handler, nil
}
