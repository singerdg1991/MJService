package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/service/constants"
	"github.com/hoitek/Maja-Service/internal/service/domain"
	"github.com/hoitek/Maja-Service/internal/service/models"
	"github.com/hoitek/Maja-Service/internal/service/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ServiceService struct {
	PostgresRepository ports.ServiceRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewServiceService(pDB ports.ServiceRepositoryPostgresDB, m *storage.MinIO) ServiceService {
	go minio.SetupMinIOStorage(constants.REPORT_TYPE_BUCKET_NAME, m)
	return ServiceService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ServiceService) Query(q *models.ServicesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying services", q)
	services, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ServicesQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	// Transform the response to the format that the frontend expects
	var items []*domain.Service
	for _, item := range services {
		items = append(items, &domain.Service{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *ServiceService) Create(payload *models.ServicesCreateRequestBody) (*domain.Service, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ServiceService) Delete(payload *models.ServicesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// TODO this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *ServiceService) Update(payload *models.ServicesCreateRequestBody, id int64) (*domain.Service, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ServiceService) GetServicesByIds(ids []int64) ([]*domain.Service, error) {
	return s.PostgresRepository.GetServicesByIds(ids)
}

func (s *ServiceService) FindByID(id int64) (*domain.Service, error) {
	r, err := s.Query(&models.ServicesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("report type not found")
	}
	services := r.Items.([]*domain.Service)
	return services[0], nil
}

func (s *ServiceService) FindByName(name string) (*domain.Service, error) {
	r, err := s.Query(&models.ServicesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("report type not found")
	}
	services := r.Items.([]*domain.Service)
	return services[0], nil
}
