package staff

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	aPorts "github.com/hoitek/Maja-Service/internal/ability/ports"
	aRepositories "github.com/hoitek/Maja-Service/internal/ability/repositories"
	aService "github.com/hoitek/Maja-Service/internal/ability/service"
	cPorts "github.com/hoitek/Maja-Service/internal/contracttype/ports"
	cRepositories "github.com/hoitek/Maja-Service/internal/contracttype/repositories"
	cService "github.com/hoitek/Maja-Service/internal/contracttype/service"
	lsPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	lsRepositories "github.com/hoitek/Maja-Service/internal/languageskill/repositories"
	lsService "github.com/hoitek/Maja-Service/internal/languageskill/service"
	permPorts "github.com/hoitek/Maja-Service/internal/license/ports"
	permRepositories "github.com/hoitek/Maja-Service/internal/license/repositories"
	permService "github.com/hoitek/Maja-Service/internal/license/service"
	nPorts "github.com/hoitek/Maja-Service/internal/notification/ports"
	nRepositories "github.com/hoitek/Maja-Service/internal/notification/repositories"
	nService "github.com/hoitek/Maja-Service/internal/notification/service"
	pPorts "github.com/hoitek/Maja-Service/internal/paymenttype/ports"
	pRepositories "github.com/hoitek/Maja-Service/internal/paymenttype/repositories"
	pService "github.com/hoitek/Maja-Service/internal/paymenttype/service"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	rRepositories "github.com/hoitek/Maja-Service/internal/role/repositories"
	rService "github.com/hoitek/Maja-Service/internal/role/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	sPorts "github.com/hoitek/Maja-Service/internal/section/ports"
	sRepositories "github.com/hoitek/Maja-Service/internal/section/repositories"
	sService "github.com/hoitek/Maja-Service/internal/section/service"
	stPorts "github.com/hoitek/Maja-Service/internal/shifttype/ports"
	stRepositories "github.com/hoitek/Maja-Service/internal/shifttype/repositories"
	stService "github.com/hoitek/Maja-Service/internal/shifttype/service"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"github.com/hoitek/Maja-Service/internal/staff/handlers"
	"github.com/hoitek/Maja-Service/internal/staff/ports"
	"github.com/hoitek/Maja-Service/internal/staff/repositories"
	"github.com/hoitek/Maja-Service/internal/staff/service"
	ntPorts "github.com/hoitek/Maja-Service/internal/stafftype/ports"
	ntRepositories "github.com/hoitek/Maja-Service/internal/stafftype/repositories"
	ntService "github.com/hoitek/Maja-Service/internal/stafftype/service"
	userPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the staff domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the staff domain module
var Module = &module{}

// GetStaffService returns a new instance of the staff service
func (m *module) GetStaffService(pDB *sql.DB) service.StaffService {
	// staff repository database based on the environment
	staffRepositoryPostgresDB := exp.TerIf[ports.StaffRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewStaffRepositoryStub(),
		repositories.NewStaffRepositoryPostgresDB(pDB),
	)

	// staff repository mongoDB
	staffRepositoryMongoDB := repositories.NewStaffRepositoryMongoDB(m.MongoDB)

	// staff service and inject the staff repository database and grpc
	staffService := service.NewStaffService(staffRepositoryPostgresDB, staffRepositoryMongoDB, m.MinIOStorage)

	// Return injected staff service
	return staffService
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

// GetRoleService returns a new instance of the role service
func (m *module) GetRoleService(pDB *sql.DB) rPorts.RoleService {
	// role repository database based on the environment
	roleRepositoryPostgresDB := exp.TerIf[rPorts.RoleRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		rRepositories.NewRoleRepositoryStub(),
		rRepositories.NewRoleRepositoryPostgresDB(pDB),
	)
	// role service and inject the role repository database and grpc
	roleService := exp.TerIf[rPorts.RoleService](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		rService.NewRoleServiceStub(),
		rService.NewRoleService(roleRepositoryPostgresDB, m.MinIOStorage),
	)
	return roleService
}

// GetSectionService returns a new instance of the section service
func (m *module) GetSectionService(pDB *sql.DB) sPorts.SectionService {
	// section repository database based on the environment
	sectionRepositoryPostgresDB := exp.TerIf[sPorts.SectionRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		sRepositories.NewSectionRepositoryStub(),
		sRepositories.NewSectionRepositoryPostgresDB(pDB),
	)

	// section service and inject the section repository database and grpc
	sectionService := sService.NewSectionService(sectionRepositoryPostgresDB, m.MinIOStorage)
	return &sectionService
}

// GetAbilityService returns a new instance of the ability service
func (m *module) GetAbilityService(pDB *sql.DB) aPorts.AbilityService {
	// ability repository database based on the environment
	abilityRepositoryPostgresDB := exp.TerIf[aPorts.AbilityRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		aRepositories.NewAbilityRepositoryStub(),
		aRepositories.NewAbilityRepositoryPostgresDB(pDB),
	)
	// ability service and inject the ability repository database and grpc
	abilityService := aService.NewAbilityService(abilityRepositoryPostgresDB, m.MinIOStorage)
	return &abilityService
}

// GetPaymentTypeService returns a new instance of the paymentType service
func (m *module) GetPaymentTypeService(pDB *sql.DB) pPorts.PaymentTypeService {
	// paymentType repository database based on the environment
	paymentTypeRepositoryPostgresDB := exp.TerIf[pPorts.PaymentTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		pRepositories.NewPaymentTypeRepositoryStub(),
		pRepositories.NewPaymentTypeRepositoryPostgresDB(pDB),
	)
	// paymentType service and inject the paymentType repository database and grpc
	paymentTypeService := pService.NewPaymentTypeService(paymentTypeRepositoryPostgresDB, m.MinIOStorage)
	return &paymentTypeService
}

// GetShiftTypeService returns a new instance of the ShiftType service
func (m *module) GetShiftTypeService(pDB *sql.DB) stPorts.ShiftTypeService {
	// ShiftType repository database based on the environment
	shiftTypeRepositoryPostgresDB := exp.TerIf[stPorts.ShiftTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		stRepositories.NewShiftTypeRepositoryStub(),
		stRepositories.NewShiftTypeRepositoryPostgresDB(pDB),
	)
	// ShiftType service and inject the ShiftType repository database and grpc
	shiftTypeService := stService.NewShiftTypeService(shiftTypeRepositoryPostgresDB, m.MinIOStorage)
	return &shiftTypeService
}

// GetContractTypeService returns a new instance of the ContractType service
func (m *module) GetContractTypeService(pDB *sql.DB) cPorts.ContractTypeService {
	// ContractType repository database based on the environment
	ContractTypeRepositoryPostgresDB := exp.TerIf[cPorts.ContractTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		cRepositories.NewContractTypeRepositoryStub(),
		cRepositories.NewContractTypeRepositoryPostgresDB(pDB),
	)
	// ContractType service and inject the ContractType repository database and grpc
	ContractTypeService := cService.NewContractTypeService(ContractTypeRepositoryPostgresDB, m.MinIOStorage)
	return &ContractTypeService
}

// GetLicenseService returns a new instance of the license service
func (m *module) GetLicenseService(pDB *sql.DB) permPorts.LicenseService {
	// license repository database based on the environment
	licenseRepositoryPostgresDB := exp.TerIf[permPorts.LicenseRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		permRepositories.NewLicenseRepositoryStub(),
		permRepositories.NewLicenseRepositoryPostgresDB(pDB),
	)
	// license service and inject the license repository database and grpc
	licenseService := permService.NewLicenseService(licenseRepositoryPostgresDB, m.MinIOStorage)
	return &licenseService
}

// GetStaffTypeService returns a new instance of the staffType service
func (m *module) GetStaffTypeService(pDB *sql.DB) ntPorts.StaffTypeService {
	// staffType repository database based on the environment
	staffTypeRepositoryPostgresDB := exp.TerIf[ntPorts.StaffTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		ntRepositories.NewStaffTypeRepositoryStub(),
		ntRepositories.NewStaffTypeRepositoryPostgresDB(pDB),
	)
	// staffType service and inject the staffType repository database and grpc
	staffTypeService := ntService.NewStaffTypeService(staffTypeRepositoryPostgresDB, m.MinIOStorage)
	return &staffTypeService
}

// GetLanguageSkillService returns a new instance of the languageSkill service
func (m *module) GetLanguageSkillService(pDB *sql.DB) lsPorts.LanguageSkillService {
	// languageSkill repository database based on the environment
	languageSkillRepositoryPostgresDB := exp.TerIf[lsPorts.LanguageSkillRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		lsRepositories.NewLanguageSkillRepositoryStub(),
		lsRepositories.NewLanguageSkillRepositoryPostgresDB(pDB),
	)
	// languageSkill service and inject the languageSkill repository database and grpc
	languageSkillService := lsService.NewLanguageSkillService(languageSkillRepositoryPostgresDB, m.MinIOStorage)
	return &languageSkillService
}

// GetNotificationService returns a new instance of the notification service
func (m *module) GetNotificationService(pDB *sql.DB) nPorts.NotificationService {
	// notification repository database based on the environment
	notificationRepositoryPostgresDB := exp.TerIf[nPorts.NotificationRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		nRepositories.NewNotificationRepositoryStub(),
		nRepositories.NewNotificationRepositoryPostgresDB(pDB),
	)
	// notification service and inject the notification repository database and grpc
	notificationService := nService.NewNotificationService(notificationRepositoryPostgresDB, m.MinIOStorage)
	return &notificationService
}

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.StaffConfig = &c
	return m
}

// SetDatabases sets the database for the staff domain module
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

// RegisterHTTP registers the staff domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.StaffHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.StaffHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for staff module
	r.Use(json.Middleware)

	// Create a new staff handler
	handler, err := handlers.NewStaffHandler(r,
		m.GetStaffService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
		m.GetRoleService(m.PostgresDB),
		m.GetSectionService(m.PostgresDB),
		m.GetAbilityService(m.PostgresDB),
		m.GetPaymentTypeService(m.PostgresDB),
		m.GetContractTypeService(m.PostgresDB),
		m.GetShiftTypeService(m.PostgresDB),
		m.GetLicenseService(m.PostgresDB),
		m.GetStaffTypeService(m.PostgresDB),
		m.GetLanguageSkillService(m.PostgresDB),
		m.GetNotificationService(m.PostgresDB),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.StaffHandler{}, err
	}

	// Return the staff handler
	return handler, nil
}

func (m *module) RegisterWorkers() *module {
	return m
}

func (m *module) ForceMigrateAndSeed() *module {
	m.Config.ForceMigrateAndSeed = true
	config.StaffConfig.ForceMigrateAndSeed = true
	return m
}
