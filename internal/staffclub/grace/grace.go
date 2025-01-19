package grace

import (
	"database/sql"
	"errors"

	rPorts "github.com/hoitek/Maja-Service/internal/reward/ports"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hoitek/Maja-Service/storage"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	rRepositories "github.com/hoitek/Maja-Service/internal/reward/repositories"
	rService "github.com/hoitek/Maja-Service/internal/reward/service"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/config"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/handlers"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/ports"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/repositories"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the grace domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the grace domain module
var Module = &module{}

// GetService returns a new instance of the grace service
func (m *module) GetService(pDB *sql.DB) ports.GraceService {
	// grace repository database based on the environment
	graceRepositoryPostgresDB := exp.TerIf[ports.GraceRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewGraceRepositoryStub(),
		repositories.NewGraceRepositoryPostgresDB(pDB),
	)
	// grace service and inject the grace repository database and grpc
	graceService := service.NewGraceService(graceRepositoryPostgresDB, m.MinIOStorage)
	return &graceService
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

// GetRewardService returns a new instance of the reward service
func (m *module) GetRewardService(pDB *sql.DB) rPorts.RewardService {
	// reward repository database based on the environment
	rewardRepositoryPostgresDB := exp.TerIf[rPorts.RewardRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		rRepositories.NewRewardRepositoryStub(),
		rRepositories.NewRewardRepositoryPostgresDB(pDB),
	)
	// reward service and inject the reward repository database and grpc
	rewardService := rService.NewRewardService(rewardRepositoryPostgresDB, m.MinIOStorage)
	return &rewardService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.GraceConfig = &c
	return m
}

// SetDatabases sets the databases
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

// RegisterHTTP registers the grace domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.GraceHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.GraceHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for grace module
	r.Use(json.Middleware)

	// Create a new grace handler
	handler, err := handlers.NewGraceHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetRewardService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.GraceHandler{}, err
	}

	// Return the grace handler
	return handler, nil
}
