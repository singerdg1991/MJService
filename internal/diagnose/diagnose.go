package diagnose

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/diagnose/config"
	"github.com/hoitek/Maja-Service/internal/diagnose/handlers"
	"github.com/hoitek/Maja-Service/internal/diagnose/ports"
	"github.com/hoitek/Maja-Service/internal/diagnose/repositories"
	"github.com/hoitek/Maja-Service/internal/diagnose/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the diagnose domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the diagnose domain module
var Module = &module{}

// GetService returns a new instance of the diagnose service
func (m *module) GetService(pDB *sql.DB) ports.DiagnoseService {
	// diagnose repository database based on the environment
	diagnoseRepositoryPostgresDB := exp.TerIf[ports.DiagnoseRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewDiagnoseRepositoryStub(),
		repositories.NewDiagnoseRepositoryPostgresDB(pDB),
	)
	// diagnose service and inject the diagnose repository database and grpc
	diagnoseService := service.NewDiagnoseService(diagnoseRepositoryPostgresDB, m.MinIOStorage)
	return &diagnoseService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.DiagnoseConfig = &c
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

// RegisterHTTP registers the diagnose domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.DiagnoseHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.DiagnoseHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for diagnose module
	r.Use(json.Middleware)

	// Create a new diagnose handler
	handler, err := handlers.NewDiagnoseHandler(r,
		m.GetService(m.PostgresDB),
	)
	if err != nil {
		return handlers.DiagnoseHandler{}, err
	}

	// Return the diagnose handler
	return handler, nil
}
