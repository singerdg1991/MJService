package cycle

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	csPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	csRepositories "github.com/hoitek/Maja-Service/internal/customer/repositories"
	csService "github.com/hoitek/Maja-Service/internal/customer/service"
	"github.com/hoitek/Maja-Service/internal/cycle/config"
	cycleConstants "github.com/hoitek/Maja-Service/internal/cycle/constants"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/handlers"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
	"github.com/hoitek/Maja-Service/internal/cycle/ports"
	"github.com/hoitek/Maja-Service/internal/cycle/repositories"
	"github.com/hoitek/Maja-Service/internal/cycle/service"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	rRepositories "github.com/hoitek/Maja-Service/internal/role/repositories"
	rService "github.com/hoitek/Maja-Service/internal/role/service"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	ssPorts "github.com/hoitek/Maja-Service/internal/service/ports"
	ssRepositories "github.com/hoitek/Maja-Service/internal/service/repositories"
	ssService "github.com/hoitek/Maja-Service/internal/service/service"
	sgPorts "github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	sgRepositories "github.com/hoitek/Maja-Service/internal/servicegrade/repositories"
	sgService "github.com/hoitek/Maja-Service/internal/servicegrade/service"
	stypePorts "github.com/hoitek/Maja-Service/internal/servicetype/ports"
	stypeRepositories "github.com/hoitek/Maja-Service/internal/servicetype/repositories"
	stypeService "github.com/hoitek/Maja-Service/internal/servicetype/service"
	stPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	stRepositories "github.com/hoitek/Maja-Service/internal/staff/repositories"
	stService "github.com/hoitek/Maja-Service/internal/staff/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	mbPorts "github.com/hoitek/Maja-Service/messagebroker/ports"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
)

// module is a module for the cycle domain
type module struct {
	Config        config.ConfigType
	PostgresDB    *sql.DB
	MongoDB       *mongo.Client
	MinIOStorage  *storage.MinIO
	MessageBroker mbPorts.MessageBroker
}

// Module is a global variable for the cycle domain module
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

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// GetServiceTypeService returns a new instance of the serviceType service
func (m *module) GetServiceTypeService(pDB *sql.DB) stypePorts.ServiceTypeService {
	// serviceType repository database based on the environment
	serviceTypeRepositoryPostgresDB := exp.TerIf[stypePorts.ServiceTypeRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		stypeRepositories.NewServiceTypeRepositoryStub(),
		stypeRepositories.NewServiceTypeRepositoryPostgresDB(pDB),
	)
	// serviceType service and inject the serviceType repository database and grpc
	serviceTypeService := stypeService.NewServiceTypeService(serviceTypeRepositoryPostgresDB, m.MinIOStorage)
	return &serviceTypeService
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

// GetCycleService returns a new instance of the cycle service
func (m *module) GetCycleService(pDB *sql.DB) ports.CycleService {
	// cycle repository database based on the environment
	customerService := m.GetCustomerService(pDB)
	serviceGradeService := m.GetServiceGradeService(pDB)
	cycleRepositoryPostgresDB := exp.TerIf[ports.CycleRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewCycleRepositoryStub(),
		repositories.NewCycleRepositoryPostgresDB(pDB, customerService, serviceGradeService),
	)

	// cycle service and inject the cycle repository database and grpc
	cycleService := service.NewCycleService(cycleRepositoryPostgresDB, m.MinIOStorage)
	return &cycleService
}

// GetStaffService returns a new instance of the staff service
func (m *module) GetStaffService(pDB *sql.DB) stPorts.StaffService {
	// staff repository database based on the environment
	staffRepositoryPostgresDB := exp.TerIf[stPorts.StaffRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		stRepositories.NewStaffRepositoryStub(),
		stRepositories.NewStaffRepositoryPostgresDB(pDB),
	)

	// staff repository mongoDB
	staffRepositoryMongoDB := stRepositories.NewStaffRepositoryMongoDB(m.MongoDB)

	// staff service and inject the staff repository database and grpc
	staffService := stService.NewStaffService(staffRepositoryPostgresDB, staffRepositoryMongoDB, m.MinIOStorage)

	// Return injected staff service
	return &staffService
}

// GetCustomerService returns a new instance of the customer service
func (m *module) GetCustomerService(pDB *sql.DB) csPorts.CustomerService {
	// customer repository database based on the environment
	customerRepositoryPostgresDB := exp.TerIf[csPorts.CustomerRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		csRepositories.NewCustomerRepositoryStub(),
		csRepositories.NewCustomerRepositoryPostgresDB(pDB),
	)

	// customer repository mongoDB
	customerRepositoryMongoDB := csRepositories.NewCustomerRepositoryMongoDB(m.MongoDB)

	// customer service and inject the customer repository database and grpc
	customerService := csService.NewCustomerService(customerRepositoryPostgresDB, customerRepositoryMongoDB, m.MinIOStorage)

	// Return injected customer service
	return &customerService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.CycleConfig = &c
	return m
}

// SetDatabases sets the database for the cycle domain module
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

// SetMessageBroker sets the message broker
func (m *module) SetMessageBroker(mb mbPorts.MessageBroker) *module {
	m.MessageBroker = mb
	return m
}

// RegisterQueues registers the cycle queues
func (m *module) RegisterQueues() *module {
	err := RegisterQueues(m.MessageBroker)
	if err != nil {
		panic(err)
	}
	return m
}

// RegisterHTTP registers the cycle domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.CycleHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.CycleHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for cycle module
	r.Use(json.Middleware)

	// Create a new cycle handler
	handler, err := handlers.NewCycleHandler(
		r,
		m.GetCycleService(m.PostgresDB),
		m.GetRoleService(m.PostgresDB),
		m.GetUserService(m.PostgresDB),
		m.GetStaffService(m.PostgresDB),
		m.GetCustomerService(m.PostgresDB),
		m.GetServiceService(m.PostgresDB),
		m.GetServiceTypeService(m.PostgresDB),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.CycleHandler{}, err
	}

	// Return the cycle handler
	return handler, nil
}

// RegisterCron sets up a cron job to check every minute if the current cycle has expired and if so, create a new cycle and update the next staff types
func (m *module) RegisterCron() *module {
	// Run a cron job to check every minute if the current cycle has expired and if so, create a new cycle and update the next staff types
	go func() {
		// Create a new cron job
		for {
			// Sleep for 1 minute
			time.Sleep(10 * time.Minute)

			// Get the cycle service
			cycleService := m.GetCycleService(m.PostgresDB)

			// Get the role service
			roleService := m.GetRoleService(m.PostgresDB)

			// Get the current cycle
			currentCycle, err := cycleService.GetCurrent()
			if err != nil {
				continue
			}

			// Check if the current cycle has expired
			if currentCycle != nil {
				if currentCycle.Status == cycleConstants.STATUS_EXPIRED {
					// Get the next staff types
					nextStaffTypesResponse, err := cycleService.QueryNextStaffTypes(&models.CyclesQueryNextStaffTypesRequestParams{
						CurrentCycleID: int(currentCycle.ID),
						Limit:          -1,
					})
					if err != nil {
						continue
					}
					nextStaffTypes, ok := nextStaffTypesResponse.Items.([]*domain.CycleNextStaffType)
					if !ok {
						continue
					}
					if len(nextStaffTypes) == 0 {
						continue
					}
					var (
						startDate        string
						endDate          *string
						freezePeriodDate string
						edBefore         time.Time
					)

					// Check if period length is set or cycle has an end date
					if currentCycle.EndDate != nil {
						edBefore = *currentCycle.EndDate
					}
					if currentCycle.PeriodLength != nil {
						switch *currentCycle.PeriodLength {
						case cycleConstants.CYCLE_PERIOD_LENGTH_ONE_WEEK:
							edBefore = currentCycle.StartDate.AddDate(0, 0, 7)
						case cycleConstants.CYCLE_PERIOD_LENGTH_TWO_WEEKS:
							edBefore = currentCycle.StartDate.AddDate(0, 0, 14)
						case cycleConstants.CYCLE_PERIOD_LENGTH_THREE_WEEKS:
							edBefore = currentCycle.StartDate.AddDate(0, 0, 21)
						}
					}

					// Calculate new cycle dates
					startEndDiff := edBefore.Sub(currentCycle.StartDate)
					freezeEndDiff := edBefore.Sub(currentCycle.FreezePeriodDate)
					ed := edBefore.Add(startEndDiff)
					startDate = edBefore.Format("2006-01-02")
					edp := ed.Format("2006-01-02")
					endDate = &edp
					fpd := ed.Add(-freezeEndDiff)
					freezePeriodDate = fpd.Format("2006-01-02")

					// Create a new cycle based on new data
					createdCycle, err := cycleService.Create(&models.CyclesCreateRequestBody{
						SectionID:              int(currentCycle.SectionID),
						StartDate:              startDate,
						EndDate:                endDate,
						PeriodLength:           currentCycle.PeriodLength,
						ShiftMorningStartTime:  currentCycle.ShiftMorningStartTime,
						ShiftMorningEndTime:    currentCycle.ShiftMorningEndTime,
						ShiftEveningStartTime:  currentCycle.ShiftEveningStartTime,
						ShiftEveningEndTime:    currentCycle.ShiftEveningEndTime,
						ShiftNightStartTime:    currentCycle.ShiftNightStartTime,
						ShiftNightEndTime:      currentCycle.ShiftNightEndTime,
						FreezePeriodDate:       freezePeriodDate,
						WishDays:               currentCycle.WishDays,
						StartDateAsDate:        &edBefore,
						EndDateAsDate:          &ed,
						FreezePeriodDateAsDate: &fpd,
					})
					if err != nil {
						log.Printf("-------- Error creating cycle: %v\n", err)
						continue
					}

					// Update the next staff types and next next staff types
					updateNextStaffTypes := []*models.CyclesUpdateStaffTypeRequestBody{}
					updateNextNextStaffTypes := []*models.CyclesUpdateNextStaffTypeRequestBody{}
					for _, nextStaffType := range nextStaffTypes {
						newDateTime := nextStaffType.DateTime
						newNextDateTime := nextStaffType.DateTime.Add(startEndDiff)
						role := roleService.GetRoleByID(int(nextStaffType.RoleID))
						if role == nil {
							continue
						}
						updateNextStaffTypes = append(updateNextStaffTypes, &models.CyclesUpdateStaffTypeRequestBody{
							ID:               int(nextStaffType.ID),
							ShiftName:        nextStaffType.ShiftName,
							DateTime:         newDateTime.Format("2006-01-02"),
							StartHour:        nextStaffType.StartHour.Format("15:04"),
							EndHour:          nextStaffType.EndHour.Format("15:04"),
							NeededStaffCount: int(nextStaffType.NeededStaffCount),
							RoleID:           int(nextStaffType.RoleID),
							DateTimeAsDate:   &newDateTime,
							Role: &domain.CycleStaffTypeRole{
								ID:   role.ID,
								Name: role.Name,
							},
							StartHourAsTime: &nextStaffType.StartHour,
							EndHourAsTime:   &nextStaffType.EndHour,
						})
						updateNextNextStaffTypes = append(updateNextNextStaffTypes, &models.CyclesUpdateNextStaffTypeRequestBody{
							ShiftName:        nextStaffType.ShiftName,
							DateTime:         newNextDateTime.Format("2006-01-02"),
							StartHour:        nextStaffType.StartHour.Format("15:04"),
							EndHour:          nextStaffType.EndHour.Format("15:04"),
							NeededStaffCount: int(nextStaffType.NeededStaffCount),
							RoleID:           int(nextStaffType.RoleID),
							DateTimeAsDate:   &newNextDateTime,
							Role: &domain.CycleNextStaffTypeRole{
								ID:   role.ID,
								Name: role.Name,
							},
							StartHourAsTime: &nextStaffType.StartHour,
							EndHourAsTime:   &nextStaffType.EndHour,
						})
					}

					// Update the new cycle next staff types
					if len(updateNextNextStaffTypes) > 0 {
						_, err := cycleService.UpdateNextStaffTypes(&models.CyclesUpdateNextStaffTypesRequestBody{
							StaffTypes: updateNextNextStaffTypes,
						}, int64(createdCycle.ID))
						if err != nil {
							log.Printf("-------- Error updating next staff types: %v\n", err)
							continue
						}
					}
					if len(updateNextStaffTypes) > 0 {
						for _, nextStaffType := range updateNextStaffTypes {
							_, err := cycleService.UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(nextStaffType, int64(createdCycle.ID), int64(currentCycle.ID))
							if err != nil {
								log.Printf("-------- Error updating staff types: %v\n", err)
								continue
							}
						}
					}
					continue
				}
			}
		}
	}()
	return m
}
