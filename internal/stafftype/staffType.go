package stafftype

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/stafftype/config"
	"github.com/hoitek/Maja-Service/internal/stafftype/handlers"
	"github.com/hoitek/Maja-Service/internal/stafftype/ports"
	"github.com/hoitek/Maja-Service/internal/stafftype/repositories"
	"github.com/hoitek/Maja-Service/internal/stafftype/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the staffType domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the staffType domain module
var Module = &module{}

// GetService returns a new instance of the staffType service
func (m *module) GetService(pDB *sql.DB) ports.StaffTypeService {
	// staffType repository database based on the environment
	staffTypeRepositoryPostgresDB := exp.TerIf[ports.StaffTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewStaffTypeRepositoryStub(),
		repositories.NewStaffTypeRepositoryPostgresDB(pDB),
	)
	// staffType service and inject the staffType repository database and grpc
	staffTypeService := service.NewStaffTypeService(staffTypeRepositoryPostgresDB, m.MinIOStorage)
	return &staffTypeService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.StaffTypeConfig = &c
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

// RegisterHTTP registers the staffType domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.StaffTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.StaffTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for staffType module
	r.Use(json.Middleware)

	// Create a new staffType handler
	handler, err := handlers.NewStaffTypeHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.StaffTypeHandler{}, err
	}

	// Return the staffType handler
	return handler, nil
}
