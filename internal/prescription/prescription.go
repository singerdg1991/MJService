package prescription

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/prescription/config"
	"github.com/hoitek/Maja-Service/internal/prescription/handlers"
	"github.com/hoitek/Maja-Service/internal/prescription/ports"
	"github.com/hoitek/Maja-Service/internal/prescription/repositories"
	"github.com/hoitek/Maja-Service/internal/prescription/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the prescription domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the prescription domain module
var Module = &module{}

// GetService returns a new instance of the prescription service
func (m *module) GetService(pDB *sql.DB) ports.PrescriptionService {
	// prescription repository database based on the environment
	prescriptionRepositoryPostgresDB := exp.TerIf[ports.PrescriptionRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewPrescriptionRepositoryStub(),
		repositories.NewPrescriptionRepositoryPostgresDB(pDB),
	)
	// prescription service and inject the prescription repository database and grpc
	prescriptionService := service.NewPrescriptionService(prescriptionRepositoryPostgresDB, m.MinIOStorage)
	return &prescriptionService
}

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.PrescriptionConfig = &c
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

// RegisterHTTP registers the prescription domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.PrescriptionHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.PrescriptionHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for prescription module
	r.Use(json.Middleware)

	// Create a new prescription handler
	handler, err := handlers.NewPrescriptionHandler(r,
		m.GetService(m.PostgresDB),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.PrescriptionHandler{}, err
	}

	// Return the prescription handler
	return handler, nil
}
