package archive

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/archive/config"
	"github.com/hoitek/Maja-Service/internal/archive/handlers"
	"github.com/hoitek/Maja-Service/internal/archive/ports"
	"github.com/hoitek/Maja-Service/internal/archive/repositories"
	"github.com/hoitek/Maja-Service/internal/archive/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	userPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the archive domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the archive domain module
var Module = &module{}

// GetService returns a new instance of the archive service
func (m *module) GetService(pDB *sql.DB) ports.ArchiveService {
	// archive repository database based on the environment
	archiveRepositoryPostgresDB := exp.TerIf[ports.ArchiveRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewArchiveRepositoryStub(),
		repositories.NewArchiveRepositoryPostgresDB(pDB),
	)
	// archive service and inject the archive repository database and grpc
	archiveService := service.NewArchiveService(archiveRepositoryPostgresDB, m.MinIOStorage)
	return &archiveService
}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB) userPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[userPorts.UserRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		uRepositories.NewUserRepositoryStub(),
		uRepositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := uRepositories.NewUserRepositoryMongoDB(m.MongoDB)

	// user service and inject the user repository database and grpc
	userService := uService.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, m.MinIOStorage)
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
	config.ArchiveConfig = &c
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// SetDatabases sets the database for the user domain module
func (m *module) SetDatabases(pDB *sql.DB, mDB *mongo.Client) *module {
	m.PostgresDB = pDB
	m.MongoDB = mDB
	return m
}

// RegisterHTTP registers the archive domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.ArchiveHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.ArchiveHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for archive module
	r.Use(json.Middleware)

	// Create a new archive handler
	handler, err := handlers.NewArchiveHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.ArchiveHandler{}, err
	}

	// Return the archive handler
	return handler, nil
}
