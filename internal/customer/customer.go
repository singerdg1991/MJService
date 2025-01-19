package customer

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	aPorts "github.com/hoitek/Maja-Service/internal/ability/ports"
	aRepositories "github.com/hoitek/Maja-Service/internal/ability/repositories"
	aService "github.com/hoitek/Maja-Service/internal/ability/service"
	adPorts "github.com/hoitek/Maja-Service/internal/address/ports"
	adRepositories "github.com/hoitek/Maja-Service/internal/address/repositories"
	adService "github.com/hoitek/Maja-Service/internal/address/service"
	cPorts "github.com/hoitek/Maja-Service/internal/contracttype/ports"
	cRepositories "github.com/hoitek/Maja-Service/internal/contracttype/repositories"
	cService "github.com/hoitek/Maja-Service/internal/contracttype/service"
	"github.com/hoitek/Maja-Service/internal/customer/config"
	"github.com/hoitek/Maja-Service/internal/customer/handlers"
	"github.com/hoitek/Maja-Service/internal/customer/ports"
	"github.com/hoitek/Maja-Service/internal/customer/repositories"
	"github.com/hoitek/Maja-Service/internal/customer/service"
	lssPorts "github.com/hoitek/Maja-Service/internal/languageskill/ports"
	lsRepositories "github.com/hoitek/Maja-Service/internal/languageskill/repositories"
	lsService "github.com/hoitek/Maja-Service/internal/languageskill/service"
	lPorts "github.com/hoitek/Maja-Service/internal/limitation/ports"
	lRepositories "github.com/hoitek/Maja-Service/internal/limitation/repositories"
	lService "github.com/hoitek/Maja-Service/internal/limitation/service"
	pPorts "github.com/hoitek/Maja-Service/internal/paymenttype/ports"
	pRepositories "github.com/hoitek/Maja-Service/internal/paymenttype/repositories"
	pService "github.com/hoitek/Maja-Service/internal/paymenttype/service"
	permPorts "github.com/hoitek/Maja-Service/internal/permission/ports"
	permRepositories "github.com/hoitek/Maja-Service/internal/permission/repositories"
	permService "github.com/hoitek/Maja-Service/internal/permission/service"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	rRepositories "github.com/hoitek/Maja-Service/internal/role/repositories"
	rService "github.com/hoitek/Maja-Service/internal/role/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	sPorts "github.com/hoitek/Maja-Service/internal/section/ports"
	sRepositories "github.com/hoitek/Maja-Service/internal/section/repositories"
	sService "github.com/hoitek/Maja-Service/internal/section/service"
	ssPorts "github.com/hoitek/Maja-Service/internal/service/ports"
	ssRepositories "github.com/hoitek/Maja-Service/internal/service/repositories"
	ssService "github.com/hoitek/Maja-Service/internal/service/service"
	sgsPorts "github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	sgsRepositories "github.com/hoitek/Maja-Service/internal/servicegrade/repositories"
	sgsService "github.com/hoitek/Maja-Service/internal/servicegrade/service"
	stPorts "github.com/hoitek/Maja-Service/internal/shifttype/ports"
	stRepositories "github.com/hoitek/Maja-Service/internal/shifttype/repositories"
	stService "github.com/hoitek/Maja-Service/internal/shifttype/service"
	stfPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	stfRepositories "github.com/hoitek/Maja-Service/internal/staff/repositories"
	stfService "github.com/hoitek/Maja-Service/internal/staff/service"
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

// module is a module for the customer domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MongoDB      *mongo.Client
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the customer domain module
var Module = &module{}

// GetCustomerService returns a new instance of the customer service
func (m *module) GetCustomerService(pDB *sql.DB) ports.CustomerService {
	// customer repository database based on the environment
	customerRepositoryPostgresDB := exp.TerIf[ports.CustomerRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewCustomerRepositoryStub(),
		repositories.NewCustomerRepositoryPostgresDB(pDB),
	)

	// customer repository mongoDB
	customerRepositoryMongoDB := repositories.NewCustomerRepositoryMongoDB(m.MongoDB)

	// customer service and inject the customer repository database and grpc
	customerService := service.NewCustomerService(customerRepositoryPostgresDB, customerRepositoryMongoDB, m.MinIOStorage)

	// Return injected customer service
	return &customerService
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

// GetLimitationService returns a new instance of the limitation service
func (m *module) GetLimitationService(pDB *sql.DB) lPorts.LimitationService {
	// limitation repository database based on the environment
	limitationRepositoryPostgresDB := exp.TerIf[lPorts.LimitationRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		lRepositories.NewLimitationRepositoryStub(),
		lRepositories.NewLimitationRepositoryPostgresDB(pDB),
	)
	// limitation service and inject the limitation repository database and grpc
	limitationService := lService.NewLimitationService(limitationRepositoryPostgresDB, m.MinIOStorage)
	return &limitationService
}

// GetAddressService returns a new instance of the address service
func (m *module) GetAddressService(pDB *sql.DB) adPorts.AddressService {
	// address repository database based on the environment
	addressRepositoryPostgresDB := exp.TerIf[adPorts.AddressRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		adRepositories.NewAddressRepositoryStub(),
		adRepositories.NewAddressRepositoryPostgresDB(pDB),
	)
	// address service and inject the address repository database and grpc
	addressService := adService.NewAddressService(addressRepositoryPostgresDB, m.MinIOStorage)
	return &addressService
}

// GetServiceService returns a new instance of the service service
func (m *module) GetServiceService(pDB *sql.DB) ssPorts.ServiceService {
	// service repository database based on the environment
	serviceRepositoryPostgresDB := exp.TerIf[ssPorts.ServiceRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		ssRepositories.NewServiceRepositoryStub(),
		ssRepositories.NewServiceRepositoryPostgresDB(pDB),
	)
	// service service and inject the service repository database and grpc
	serviceService := ssService.NewServiceService(serviceRepositoryPostgresDB, m.MinIOStorage)
	return &serviceService
}

// GetServiceGradeService returns a new instance of the serviceGrade service
func (m *module) GetServiceGradeService(pDB *sql.DB) sgsPorts.ServiceGradeService {
	// serviceGrade repository database based on the environment
	serviceGradeRepositoryPostgresDB := exp.TerIf[sgsPorts.ServiceGradeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		sgsRepositories.NewServiceGradeRepositoryStub(),
		sgsRepositories.NewServiceGradeRepositoryPostgresDB(pDB),
	)
	// serviceGrade service and inject the serviceGrade repository database and grpc
	serviceGradeService := sgsService.NewServiceGradeService(serviceGradeRepositoryPostgresDB, m.MinIOStorage)
	return &serviceGradeService
}

// GetStaffService returns a new instance of the staff service
func (m *module) GetStaffService(pDB *sql.DB) stfPorts.StaffService {
	// staff repository database based on the environment
	staffRepositoryPostgresDB := exp.TerIf[stfPorts.StaffRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		stfRepositories.NewStaffRepositoryStub(),
		stfRepositories.NewStaffRepositoryPostgresDB(pDB),
	)

	// staff repository mongoDB
	staffRepositoryMongoDB := stfRepositories.NewStaffRepositoryMongoDB(m.MongoDB)

	// staff service and inject the staff repository database and grpc
	staffService := stfService.NewStaffService(staffRepositoryPostgresDB, staffRepositoryMongoDB, m.MinIOStorage)

	// Return injected staff service
	return &staffService
}

// GetLanguageSkillService returns a new instance of the languageskill service
func (m *module) GetLanguageSkillService(pDB *sql.DB) lssPorts.LanguageSkillService {
	// languageskill repository database based on the environment
	languageskillRepositoryPostgresDB := exp.TerIf[lssPorts.LanguageSkillRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		lsRepositories.NewLanguageSkillRepositoryStub(),
		lsRepositories.NewLanguageSkillRepositoryPostgresDB(pDB),
	)
	// languageskill service and inject the languageskill repository database and grpc
	languageskillService := lsService.NewLanguageSkillService(languageskillRepositoryPostgresDB, m.MinIOStorage)
	return &languageskillService
}

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.CustomerConfig = &c
	return m
}

// SetDatabases sets the database for the customer domain module
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

// RegisterHTTP registers the customer domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.CustomerHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.CustomerHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for customer module
	r.Use(json.Middleware)

	// Create a new customer handler
	handler, err := handlers.NewCustomerHandler(r,
		m.GetCustomerService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
		m.GetRoleService(m.PostgresDB),
		m.GetSectionService(m.PostgresDB),
		m.GetLimitationService(m.PostgresDB),
		m.GetAddressService(m.PostgresDB),
		m.GetServiceService(m.PostgresDB),
		m.GetStaffService(m.PostgresDB),
		m.GetServiceGradeService(m.PostgresDB),
		m.GetLanguageSkillService(m.PostgresDB),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.CustomerHandler{}, err
	}

	// Return the customer handler
	return handler, nil
}

func (m *module) RegisterWorkers() *module {
	return m
}

func (m *module) ForceMigrateAndSeed() *module {
	m.Config.ForceMigrateAndSeed = true
	config.CustomerConfig.ForceMigrateAndSeed = true
	return m
}
