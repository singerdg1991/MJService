package role

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	permPorts "github.com/hoitek/Maja-Service/internal/permission/ports"
	permRepositories "github.com/hoitek/Maja-Service/internal/permission/repositories"
	permService "github.com/hoitek/Maja-Service/internal/permission/service"
	"github.com/hoitek/Maja-Service/internal/role/config"
	"github.com/hoitek/Maja-Service/internal/role/handlers"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	"github.com/hoitek/Maja-Service/internal/role/repositories"
	"github.com/hoitek/Maja-Service/internal/role/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
)

// module is a module for the role domain
type module struct {
	Config       config.ConfigType
	PostgresDB   *sql.DB
	MinIOStorage *storage.MinIO
}

// Module is a global variable for the role domain module
var Module = &module{}

// GetService returns a new instance of the role service
func (m *module) GetService(pDB *sql.DB) rPorts.RoleService {
	// role repository database based on the environment
	roleRepositoryPostgresDB := exp.TerIf[rPorts.RoleRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewRoleRepositoryStub(),
		repositories.NewRoleRepositoryPostgresDB(pDB),
	)
	// role service and inject the role repository database and grpc
	roleService := exp.TerIf[rPorts.RoleService](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		service.NewRoleServiceStub(),
		service.NewRoleService(roleRepositoryPostgresDB, m.MinIOStorage),
	)
	return roleService
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
	config.RoleConfig = &c
	return m
}

// SetDatabase sets the database for the setting domain module
func (m *module) SetDatabase(pDB *sql.DB) *module {
	m.PostgresDB = pDB
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// RegisterHTTP registers the role domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.RoleHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.RoleHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for role module
	r.Use(json.Middleware)

	// Create a new role handler
	handler, err := handlers.NewRoleHandler(
		r,
		m.GetService(m.PostgresDB),
		m.GetPermissionService(m.PostgresDB),
	)
	if err != nil {
		return handlers.RoleHandler{}, err
	}

	// Return the role handler
	return handler, nil
}
