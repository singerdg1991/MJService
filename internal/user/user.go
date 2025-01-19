package user

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	lPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	lRepositories "github.com/hoitek/Maja-Service/internal/languageskill/repositories"
	lService "github.com/hoitek/Maja-Service/internal/languageskill/service"
	oPorts "github.com/hoitek/Maja-Service/internal/otp/ports"
	oRepositories "github.com/hoitek/Maja-Service/internal/otp/repositories"
	oService "github.com/hoitek/Maja-Service/internal/otp/service"
	rolePorts "github.com/hoitek/Maja-Service/internal/role/ports"
	repoRole "github.com/hoitek/Maja-Service/internal/role/repositories"
	roleService "github.com/hoitek/Maja-Service/internal/role/service"
	nPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	nRepositories "github.com/hoitek/Maja-Service/internal/staff/repositories"
	nService "github.com/hoitek/Maja-Service/internal/staff/service"
	"github.com/hoitek/Maja-Service/internal/user/config"
	"github.com/hoitek/Maja-Service/internal/user/handlers"
	"github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/internal/user/repositories"
	"github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the user domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the user domain module
var Module = &module{}

// GetUserService returns a new instance of the user service
func (m *module) GetUserService(pDB *sql.DB) ports.UserService {
	// user repository database based on the environment
	userRepositoryPostgresDB := exp.TerIf[ports.UserRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewUserRepositoryStub(),
		repositories.NewUserRepositoryPostgresDB(pDB),
	)

	// user repository mongoDB
	userRepositoryMongoDB := repositories.NewUserRepositoryMongoDB(m.MongoDB)

	// user service and inject the user repository database and grpc
	userService := service.NewUserService(userRepositoryPostgresDB, userRepositoryMongoDB, m.MinIOStorage)
	return userService
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
	return &staffService
}

// GetLanguageSkillService returns a new instance of the languageskill service
func (m *module) GetLanguageSkillService(pDB *sql.DB) lPorts.LanguageSkillService {
	// languageskill repository database based on the environment
	languageskillRepositoryPostgresDB := exp.TerIf[lPorts.LanguageSkillRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		lRepositories.NewLanguageSkillRepositoryStub(),
		lRepositories.NewLanguageSkillRepositoryPostgresDB(pDB),
	)
	// languageskill service and inject the languageskill repository database and grpc
	languageskillService := lService.NewLanguageSkillService(languageskillRepositoryPostgresDB, m.MinIOStorage)
	return &languageskillService
}

func (m *module) GetRoleService(pDB *sql.DB) rolePorts.RoleService {
	// role repository database based on the environment
	roleRepositoryPostgresDB := exp.TerIf[rolePorts.RoleRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repoRole.NewRoleRepositoryStub(),
		repoRole.NewRoleRepositoryPostgresDB(pDB),
	)

	// user service and inject the user repository database and grpc
	roleService := exp.TerIf[rolePorts.RoleService](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		roleService.NewRoleServiceStub(),
		roleService.NewRoleService(roleRepositoryPostgresDB, m.MinIOStorage),
	)
	return roleService
}

// GetOTPService returns a new instance of the otp service
func (m *module) GetOTPService(pDB *sql.DB) oPorts.OTPService {
	// otp repository database based on the environment
	otpRepositoryPostgresDB := exp.TerIf[oPorts.OTPRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		oRepositories.NewOTPRepositoryStub(),
		oRepositories.NewOTPRepositoryPostgresDB(pDB),
	)
	// otp service and inject the otp repository database and grpc
	otpService := oService.NewOTPService(otpRepositoryPostgresDB)
	return otpService
}

// Setup sets up the user domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.UserConfig = &c
	return m
}

// SetDatabases sets the database for the user domain module
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

// RegisterHTTP registers the user domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.UserHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.UserHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for user module
	r.Use(json.Middleware)

	// Create a new user handler
	handler, err := handlers.NewUserHandler(
		r,
		m.GetUserService(m.PostgresDB),
		m.GetRoleService(m.PostgresDB),
		m.GetLanguageSkillService(m.PostgresDB),
		m.GetStaffService(m.PostgresDB),
		m.GetOTPService(m.PostgresDB),
	)
	if err != nil {
		return handlers.UserHandler{}, err
	}

	// Return the user handler
	return handler, nil
}

func (m *module) RegisterWorkers() *module {
	return m
}
