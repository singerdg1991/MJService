package ability

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/ability/config"
	"github.com/hoitek/Maja-Service/internal/ability/handlers"
	"github.com/hoitek/Maja-Service/internal/ability/ports"
	"github.com/hoitek/Maja-Service/internal/ability/repositories"
	"github.com/hoitek/Maja-Service/internal/ability/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the ability domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the ability domain module
var Module = &module{}

// GetService returns a new instance of the ability service
func (m *module) GetService(pDB *sql.DB) service.AbilityService {
	// ability repository database based on the environment
	abilityRepositoryPostgresDB := exp.TerIf[ports.AbilityRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewAbilityRepositoryStub(),
		repositories.NewAbilityRepositoryPostgresDB(pDB),
	)
	// ability service and inject the ability repository database and grpc
	abilityService := service.NewAbilityService(abilityRepositoryPostgresDB, m.MinIOStorage)
	return abilityService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.AbilityConfig = &c
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

// RegisterHTTP registers the ability domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.AbilityHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.AbilityHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for ability module
	r.Use(json.Middleware)

	// Create a new ability handler
	handler, err := handlers.NewAbilityHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.AbilityHandler{}, err
	}

	// Return the ability handler
	return handler, nil
}
