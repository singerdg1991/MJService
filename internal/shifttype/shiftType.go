package shifttype

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/shifttype/config"
	"github.com/hoitek/Maja-Service/internal/shifttype/handlers"
	"github.com/hoitek/Maja-Service/internal/shifttype/ports"
	"github.com/hoitek/Maja-Service/internal/shifttype/repositories"
	"github.com/hoitek/Maja-Service/internal/shifttype/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the ShiftType domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the ShiftType domain module
var Module = &module{}

// GetService returns a new instance of the ShiftType service
func (m *module) GetService(pDB *sql.DB) service.ShiftTypeService {
	// ShiftType repository database based on the environment
	shiftTypeRepositoryPostgresDB := exp.TerIf[ports.ShiftTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewShiftTypeRepositoryStub(),
		repositories.NewShiftTypeRepositoryPostgresDB(pDB),
	)
	// ShiftType service and inject the ShiftType repository database and grpc
	shiftTypeService := service.NewShiftTypeService(shiftTypeRepositoryPostgresDB, m.MinIOStorage)
	return shiftTypeService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.ShiftTypeConfig = &c
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

// RegisterHTTP registers the ShiftType domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ShiftTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ShiftTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for ShiftType module
	r.Use(json.Middleware)

	// Create a new ShiftType handler
	handler, err := handlers.NewShiftTypeHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.ShiftTypeHandler{}, err
	}

	// Return the ShiftType handler
	return handler, nil
}
