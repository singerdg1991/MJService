package paymenttype

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/paymenttype/config"
	"github.com/hoitek/Maja-Service/internal/paymenttype/handlers"
	"github.com/hoitek/Maja-Service/internal/paymenttype/ports"
	"github.com/hoitek/Maja-Service/internal/paymenttype/repositories"
	"github.com/hoitek/Maja-Service/internal/paymenttype/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the paymentType domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the paymentType domain module
var Module = &module{}

// GetService returns a new instance of the paymentType service
func (m *module) GetService(pDB *sql.DB) service.PaymentTypeService {
	// paymentType repository database based on the environment
	paymentTypeRepositoryPostgresDB := exp.TerIf[ports.PaymentTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewPaymentTypeRepositoryStub(),
		repositories.NewPaymentTypeRepositoryPostgresDB(pDB),
	)
	// paymentType service and inject the paymentType repository database and grpc
	paymentTypeService := service.NewPaymentTypeService(paymentTypeRepositoryPostgresDB, m.MinIOStorage)
	return paymentTypeService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.PaymentTypeConfig = &c
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

// RegisterHTTP registers the paymentType domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.PaymentTypeHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.PaymentTypeHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for paymentType module
	r.Use(json.Middleware)

	// Create a new paymentType handler
	handler, err := handlers.NewPaymentTypeHandler(r, m.GetService(m.PostgresDB))
	if err != nil {
		return handlers.PaymentTypeHandler{}, err
	}

	// Return the paymentType handler
	return handler, nil
}
