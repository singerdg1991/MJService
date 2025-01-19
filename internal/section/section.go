package section

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/constants"
	"github.com/hoitek/Maja-Service/internal/section/config"
	"github.com/hoitek/Maja-Service/internal/section/handlers"
	"github.com/hoitek/Maja-Service/internal/section/ports"
	"github.com/hoitek/Maja-Service/internal/section/repositories"
	"github.com/hoitek/Maja-Service/internal/section/service"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Middlewares/json"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// module is a module for the section domain
type module struct {
	Config         config.ConfigType
	PostgresDB     *sql.DB
	MongoDB        *mongo.Client
	GRPCConnection *grpc.ClientConn
	MinIOStorage   *storage.MinIO
}

// Module is a global variable for the section domain module
var Module = &module{}

// GetService returns a new instance of the section service
func (m *module) GetService(pDB *sql.DB, gConn *grpc.ClientConn) service.SectionService {
	// section repository database based on the environment
	sectionRepositoryPostgresDB := exp.TerIf[ports.SectionRepositoryPostgresDB](
		m.Config.Environment == constants.ENVIRONMENT_TESTING,
		repositories.NewSectionRepositoryStub(),
		repositories.NewSectionRepositoryPostgresDB(pDB),
	)

	// section service and inject the section repository database and grpc
	sectionService := service.NewSectionService(sectionRepositoryPostgresDB, m.MinIOStorage)
	return sectionService
}

// Setup sets up the setting domain module
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.SectionConfig = &c
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

// SetGRPCConnection sets the grpc connection for the staff domain module
func (m *module) SetGRPCConnection(conn *grpc.ClientConn) *module {
	m.GRPCConnection = conn
	return m
}

// RegisterHTTP registers the section domain http routes
func (m *module) RegisterHTTP(r *mux.Router) (handlers.SectionHandler, error) {
	// Check if router is nil
	if r == nil {
		return handlers.SectionHandler{}, errors.New("router can not be nil")
	}

	// Create a new group router to scope any middleware here and prevent polluting the global router
	r = r.PathPrefix("/").Subrouter()

	// Use Json middleware for section module
	r.Use(json.Middleware)

	// Create a new section handler
	handler, err := handlers.NewSectionHandler(r, m.GetService(m.PostgresDB, m.GRPCConnection))
	if err != nil {
		return handlers.SectionHandler{}, err
	}

	// Return the section handler
	return handler, nil
}
