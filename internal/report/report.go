package report

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	cPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	cRepo "github.com/hoitek/Maja-Service/internal/customer/repositories"
	cService "github.com/hoitek/Maja-Service/internal/customer/service"
	cycleRepo "github.com/hoitek/Maja-Service/internal/cycle/repositories"
	"github.com/hoitek/Maja-Service/internal/report/config"
	"github.com/hoitek/Maja-Service/internal/report/handlers"
	rPorts "github.com/hoitek/Maja-Service/internal/report/ports"
	rRepo "github.com/hoitek/Maja-Service/internal/report/repositories"
	rService "github.com/hoitek/Maja-Service/internal/report/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	sgPorts "github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	sgRepositories "github.com/hoitek/Maja-Service/internal/servicegrade/repositories"
	sgService "github.com/hoitek/Maja-Service/internal/servicegrade/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the report domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the report domain module
var Module = &module{}

// GetCustomerService returns a new instance of the customer service
func (m *module) GetCustomerService(pDB *sql.DB) cPorts.CustomerService {
	// customer repository database based on the environment
	customerRepositoryPostgresDB := exp.TerIf[cPorts.CustomerRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		cRepo.NewCustomerRepositoryStub(),
		cRepo.NewCustomerRepositoryPostgresDB(pDB),
	)

	// customer repository mongoDB
	customerRepositoryMongoDB := cRepo.NewCustomerRepositoryMongoDB(m.MongoDB)

	// customer service and inject the customer repository database and grpc
	customerService := cService.NewCustomerService(customerRepositoryPostgresDB, customerRepositoryMongoDB, m.MinIOStorage)

	// Return injected customer service
	return &customerService
}

// GetServiceGradeService returns a new instance of the servicegrade service
func (m *module) GetServiceGradeService(pDB *sql.DB) sgPorts.ServiceGradeService {
	// servicegrade repository database based on the environment
	servicegradeRepositoryPostgresDB := exp.TerIf[sgPorts.ServiceGradeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		sgRepositories.NewServiceGradeRepositoryStub(),
		sgRepositories.NewServiceGradeRepositoryPostgresDB(pDB),
	)
	// servicegrade service and inject the servicegrade repository database and grpc
	servicegradeService := sgService.NewServiceGradeService(servicegradeRepositoryPostgresDB, m.MinIOStorage)
	return &servicegradeService
}

// GetReportService returns a new instance of the report service
func (m *module) GetReportService(pDB *sql.DB) rPorts.ReportService {
	// cycle repository dependencies
	customerService := m.GetCustomerService(pDB)
	serviceGradeService := m.GetServiceGradeService(pDB)
	cycleRepository := cycleRepo.NewCycleRepositoryPostgresDB(pDB, customerService, serviceGradeService)

	// report repository database based on the environment
	reportRepositoryPostgresDB := exp.TerIf[rPorts.ReportRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		rRepo.NewReportRepositoryStub(),
		rRepo.NewReportRepositoryPostgresDB(pDB, cycleRepository),
	)

	// report service and inject the report repository database and grpc
	reportService := rService.NewReportService(reportRepositoryPostgresDB, m.MinIOStorage)

	// Return injected report service
	return &reportService
}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB, mDB *mongo.Client) uPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[uPorts.UserRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		uRepositories.NewUserRepositoryStub(),
		uRepositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := uRepositories.NewUserRepositoryMongoDB(mDB)

	// user service and inject the user repository database and grpc
	userService := uService.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, storage.MinIOStorage)
	return userService
}

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.ReportConfig = &c
	return m
}

// SetDatabases sets the databases for the report domain module
func (m *module) SetDatabases(pDB *sql.DB, mDB *mongo.Client) *module {
	m.PostgresDB = pDB
	m.MongoDB = mDB
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// RegisterHTTP registers the report domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ReportHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ReportHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for report module
	r.Use(json.Middleware)

	// Create a new report handler
	handler, err := handlers.NewReportHandler(r,
		m.GetCustomerService(m.PostgresDB),
		m.GetReportService(m.PostgresDB),
		m.GetS3Service(),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.ReportHandler{}, err
	}

	// Return the report handler
	return handler, nil
}

func (m *module) RegisterWorkers() *module {
	return m
}

func (m *module) ForceMigrateAndSeed() *module {
	m.Config.ForceMigrateAndSeed = true
	config.ReportConfig.ForceMigrateAndSeed = true
	return m
}
