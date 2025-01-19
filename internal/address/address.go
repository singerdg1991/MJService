package address

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/address/config"
	"github.com/hoitek/Maja-Service/internal/address/handlers"
	"github.com/hoitek/Maja-Service/internal/address/ports"
	"github.com/hoitek/Maja-Service/internal/address/repositories"
	"github.com/hoitek/Maja-Service/internal/address/service"
	cPorts "github.com/hoitek/Maja-Service/internal/city/ports"
	cRepositories "github.com/hoitek/Maja-Service/internal/city/repositories"
	cService "github.com/hoitek/Maja-Service/internal/city/service"
	csPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	csRepositories "github.com/hoitek/Maja-Service/internal/customer/repositories"
	csService "github.com/hoitek/Maja-Service/internal/customer/service"
	nPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	nRepositories "github.com/hoitek/Maja-Service/internal/staff/repositories"
	nService "github.com/hoitek/Maja-Service/internal/staff/service"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	uRepositories "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// module is a module for the address domain
type module struct {
	Config         config.ConfigType
	PostgresDB     *sql.DB
	MongoDB        *mongo.Client
	GRPCConnection *grpc.ClientConn
	MinIOStorage   *storage.MinIO
}

// Module is a global variable for the address domain module
var Module = &module{}

// GetService returns a new instance of the address service
func (m *module) GetService(pDB *sql.DB) service.AddressService {
	// address repository database based on the environment
	addressRepositoryPostgresDB := exp.TerIf[ports.AddressRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewAddressRepositoryStub(),
		repositories.NewAddressRepositoryPostgresDB(pDB),
	)
	// address service and inject the address repository database and grpc
	addressService := service.NewAddressService(addressRepositoryPostgresDB, m.MinIOStorage)
	return addressService
}

// GetStaffService returns a new instance of the staff service
func (m *module) GetStaffService(pDB *sql.DB, gConn *grpc.ClientConn) nPorts.StaffService {
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

// GetCityService returns a new instance of the city service
func (m *module) GetCityService(pDB *sql.DB) cPorts.CityService {
	// city repository database based on the environment
	cityRepositoryPostgresDB := exp.TerIf[cPorts.CityRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		cRepositories.NewCityRepositoryStub(),
		cRepositories.NewCityRepositoryPostgresDB(pDB),
	)
	// city service and inject the city repository database and grpc
	cityService := cService.NewCityService(cityRepositoryPostgresDB, m.MinIOStorage)
	return &cityService
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
	config.AddressConfig = &c
	return m
}

// SetGRPCConnection sets the grpc connection for the staff domain module
func (m *module) SetGRPCConnection(conn *grpc.ClientConn) *module {
	m.GRPCConnection = conn
	return m
}

// SetDatabases sets the database for the address domain module
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

// RegisterHTTP registers the address domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.AddressHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.AddressHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for address module
	r.Use(json.Middleware)

	// Create a new address handler
	handler, err := handlers.NewAddressHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetCityService(m.PostgresDB),
		m.GetStaffService(m.PostgresDB, m.GRPCConnection),
		m.GetCustomerService(m.PostgresDB),
		m.GetUserService(m.PostgresDB, m.MongoDB),
	)
	if err != nil {
		return handlers.AddressHandler{}, err
	}

	// Return the address handler
	return handler, nil
}
