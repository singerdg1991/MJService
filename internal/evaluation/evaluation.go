package evaluation

import (
	"database/sql"
	"errors"

	nPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	nRepositories "github.com/hoitek/Maja-Service/internal/staff/repositories"
	nService "github.com/hoitek/Maja-Service/internal/staff/service"
	"github.com/hoitek/Maja-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/evaluation/config"
	"github.com/hoitek/Maja-Service/internal/evaluation/handlers"
	"github.com/hoitek/Maja-Service/internal/evaluation/ports"
	"github.com/hoitek/Maja-Service/internal/evaluation/repositories"
	"github.com/hoitek/Maja-Service/internal/evaluation/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the evaluation domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the evaluation domain module
var Module = &module{}

// GetService returns a new instance of the evaluation service
func (m *module) GetService(pDB *sql.DB) ports.EvaluationService {
	// evaluation repository database based on the environment
	evaluationRepositoryPostgresDB := exp.TerIf[ports.EvaluationRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewEvaluationRepositoryStub(),
		repositories.NewEvaluationRepositoryPostgresDB(pDB),
	)
	// evaluation service and inject the evaluation repository database and grpc
	evaluationService := service.NewEvaluationService(evaluationRepositoryPostgresDB, m.MinIOStorage)
	return &evaluationService
}

// GetStaffService returns a new instance of the staff service
func (m *module) GetStaffService(pDB *sql.DB) nPorts.StaffService {
	// staff repository database based on the environment
	staffRepositoryPostgresDB := exp.TerIf[nPorts.StaffRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		nRepositories.NewStaffRepositoryStub(),
		nRepositories.NewStaffRepositoryPostgresDB(pDB),
	)

	// staff repository mongoDB
	staffRepositoryMongoDB := nRepositories.NewStaffRepositoryMongoDB(m.MongoDB)

	// staff service and inject the staff repository database and grpc
	staffService := nService.NewStaffService(staffRepositoryPostgresDB, staffRepositoryMongoDB, m.MinIOStorage)

	// Return injected staff service
	return &staffService
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

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.EvaluationConfig = &c
	return m
}

// SetDatabases sets the database for the evaluation domain module
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

// RegisterHTTP registers the evaluation domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.EvaluationHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.EvaluationHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for evaluation module
	r.Use(json.Middleware)

	// Create a new evaluation handler
	handler, err := handlers.NewEvaluationHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetStaffService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.EvaluationHandler{}, err
	}

	// Return the evaluation handler
	return handler, nil
}
