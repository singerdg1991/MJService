package static

import (
	"github.com/gorilla/mux"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	s3Service "github.com/hoitek/Maja-Service/internal/s3/service"
	"github.com/hoitek/Maja-Service/internal/static/config"
	"github.com/hoitek/Maja-Service/internal/static/handlers"
	"github.com/hoitek/Maja-Service/internal/static/ports"
	"github.com/hoitek/Maja-Service/internal/static/service"
	"github.com/hoitek/Maja-Service/storage"
)

type module struct {
	Config       config.ConfigType
	MinIOStorage *storage.MinIO
}

var Module = &module{}

// Setup sets the config
func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.StaticConfig = &c
	return m
}

// SetMinIOStorage sets the minio storage
func (m *module) SetMinIOStorage(s *storage.MinIO) *module {
	m.MinIOStorage = s
	return m
}

// GetStaticService returns a new instance of the static service
func (m *module) GetStaticService() ports.StaticService {
	streetService := service.NewStaticService(m.MinIOStorage)
	return streetService
}

// GetS3Service returns a new instance of the s3 service
func (m *module) GetS3Service() s3Ports.S3Service {
	s3Service := s3Service.NewS3Service(m.MinIOStorage)
	return s3Service
}

// RegisterHTTP registers the http handlers
func (m *module) RegisterHTTP(r *mux.Router) (handlers.StaticHandler, error) {
	handler, err := handlers.NewStaticHandler(
		r,
		m.GetStaticService(),
		m.GetS3Service(),
	)
	if err != nil {
		return handlers.StaticHandler{}, err
	}
	return handler, nil
}
