package quiz

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	permPorts "github.com/hoitek/Maja-Service/internal/permission/ports"
	permRepositories "github.com/hoitek/Maja-Service/internal/permission/repositories"
	permService "github.com/hoitek/Maja-Service/internal/permission/service"
	"github.com/hoitek/Maja-Service/internal/quiz/config"
	"github.com/hoitek/Maja-Service/internal/quiz/handlers"
	rPorts "github.com/hoitek/Maja-Service/internal/quiz/ports"
	"github.com/hoitek/Maja-Service/internal/quiz/repositories"
	"github.com/hoitek/Maja-Service/internal/quiz/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the quiz domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the quiz domain module
var Module = &module{}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB) uPorts.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[uPorts.UserRepositoryPostgresDB](
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

// GetService returns a new instance of the quiz service
func (m *module) GetService(pDB *sql.DB) rPorts.QuizService {
	// quiz repository database based on the environment
	quizRepositoryPostgresDB := exp.TerIf[rPorts.QuizRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewQuizRepositoryStub(),
		repositories.NewQuizRepositoryPostgresDB(pDB),
	)
	// quiz service and inject the quiz repository database and grpc
	quizService := service.NewQuizService(quizRepositoryPostgresDB, m.MinIOStorage)
	return &quizService
}

// GetPermissionService returns a new instance of the permission service
func (m *module) GetPermissionService(pDB *sql.DB) permPorts.PermissionService {
	// permission repository database based on the environment
	permissionRepositoryPostgresDB := exp.TerIf[permPorts.PermissionRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		permRepositories.NewPermissionRepositoryStub(),
		permRepositories.NewPermissionRepositoryPostgresDB(pDB),
	)
	// permission service and inject the permission repository database and grpc
	permissionService := permService.NewPermissionService(permissionRepositoryPostgresDB, m.MinIOStorage)
	return &permissionService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.QuizConfig = &c
	return m
}

// SetDatabases sets the database for the quiz domain module
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

// RegisterHTTP registers the quiz domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.QuizHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.QuizHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for quiz module
	r.Use(json.Middleware)

	// Create a new quiz handler
	handler, err := handlers.NewQuizHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetPermissionService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
	)
	if err != nil {
		return handlers.QuizHandler{}, err
	}

	// Return the quiz handler
	return handler, nil
}
