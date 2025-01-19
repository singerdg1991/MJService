package city

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/city/config"
	"github.com/hoitek/Maja-Service/internal/city/handlers"
	"github.com/hoitek/Maja-Service/internal/city/ports"
	"github.com/hoitek/Maja-Service/internal/city/repositories"
	"github.com/hoitek/Maja-Service/internal/city/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the city domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the city domain module
var Module = &module{}

// GetService returns a new instance of the city service
func (m *module) GetService(pDB *sql.DB) service.CityService {
	// city repository database based on the environment
	cityRepositoryPostgresDB := exp.TerIf[ports.CityRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewCityRepositoryStub(),
		repositories.NewCityRepositoryPostgresDB(pDB),
	)
	// city service and inject the city repository database and grpc
	cityService := service.NewCityService(cityRepositoryPostgresDB, m.MinIOStorage)
	return cityService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.CityConfig = &c
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

// RegisterHTTP registers the city domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.CityHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.CityHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for city module
	r.Use(json.Middleware)

	// Create a new city handler
	handler, err := handlers.NewCityHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.CityHandler{}, err
	}

	// Return the city handler
	return handler, nil
}
