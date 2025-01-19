package geartype

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/geartype/config"
	"github.com/hoitek/Maja-Service/internal/geartype/handlers"
	"github.com/hoitek/Maja-Service/internal/geartype/ports"
	"github.com/hoitek/Maja-Service/internal/geartype/repositories"
	"github.com/hoitek/Maja-Service/internal/geartype/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the geartype domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the geartype domain module
var Module = &module{}

// GetService returns a new instance of the geartype service
func (m *module) GetService(pDB *sql.DB) service.GearTypeService {
	// geartype repository database based on the environment
	geartypeRepositoryPostgresDB := exp.TerIf[ports.GearTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewGearTypeRepositoryStub(),
		repositories.NewGearTypeRepositoryPostgresDB(pDB),
	)
	// geartype service and inject the geartype repository database and grpc
	geartypeService := service.NewGearTypeService(geartypeRepositoryPostgresDB, m.MinIOStorage)
	return geartypeService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.GearTypeConfig = &c
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// SetDatabase sets the database for the setting domain module
func (m *module) SetDatabase(pDB *sql.DB) *module {
	m.PostgresDB = pDB
	return m
}

// RegisterHTTP registers the geartype domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.GearTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.GearTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for geartype module
	r.Use(json.Middleware)

	// Create a new geartype handler
	handler, err := handlers.NewGearTypeHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.GearTypeHandler{}, err
	}

	// Return the geartype handler
	return handler, nil
}
