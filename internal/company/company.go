package company

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/company/config"
	"github.com/hoitek/Maja-Service/internal/company/handlers"
	"github.com/hoitek/Maja-Service/internal/company/ports"
	"github.com/hoitek/Maja-Service/internal/company/repositories"
	"github.com/hoitek/Maja-Service/internal/company/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the company domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the company domain module
var Module = &module{}

// GetService returns a new instance of the company service
func (m *module) GetService(pDB *sql.DB) service.CompanyService {
	// company repository database based on the environment
	companyRepositoryPostgresDB := exp.TerIf[ports.CompanyRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewCompanyRepositoryStub(),
		repositories.NewCompanyRepositoryPostgresDB(pDB),
	)
	// company service and inject the company repository database and grpc
	companyService := service.NewCompanyService(companyRepositoryPostgresDB, m.MinIOStorage)
	return companyService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.CompanyConfig = &c
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

// RegisterHTTP registers the company domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.CompanyHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.CompanyHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for company module
	r.Use(json.Middleware)

	// Create a new company handler
	handler, err := handlers.NewCompanyHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.CompanyHandler{}, err
	}

	// Return the company handler
	return handler, nil
}
