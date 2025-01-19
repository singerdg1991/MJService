package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/servicetype/constants"
	"github.com/hoitek/Maja-Service/internal/servicetype/domain"
	"github.com/hoitek/Maja-Service/internal/servicetype/models"
	"github.com/hoitek/Maja-Service/internal/servicetype/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ServiceTypeService struct {
	PostgresRepository ports.ServiceTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewServiceTypeService(pDB ports.ServiceTypeRepositoryPostgresDB, m *storage.MinIO) ServiceTypeService {
	go minio.SetupMinIOStorage(constants.REPORT_CATEGORY_BUCKET_NAME, m)
	return ServiceTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ServiceTypeService) Query(q *models.ServiceTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying serviceTypes", q)
	serviceTypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ServiceTypesQueryRequestParams{
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
	var items []*domain.ServiceType
	for _, item := range serviceTypes {
		items = append(items, &domain.ServiceType{
			ID:          item.ID,
			ServiceID:   item.ServiceID,
			Service:     item.Service,
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

func (s *ServiceTypeService) Create(payload *models.ServiceTypesCreateRequestBody) (*domain.ServiceType, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ServiceTypeService) Delete(payload *models.ServiceTypesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *ServiceTypeService) Update(payload *models.ServiceTypesCreateRequestBody, id int64) (*domain.ServiceType, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ServiceTypeService) GetServiceTypesByIds(ids []int64) ([]*domain.ServiceType, error) {
	return s.PostgresRepository.GetServiceTypesByIds(ids)
}

func (s *ServiceTypeService) FindByID(id int64) (*domain.ServiceType, error) {
	r, err := s.Query(&models.ServiceTypesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("report category not found")
	}
	serviceTypes := r.Items.([]*domain.ServiceType)
	return serviceTypes[0], nil
}

func (s *ServiceTypeService) FindByName(name string) (*domain.ServiceType, error) {
	r, err := s.Query(&models.ServiceTypesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceTypeFilterType{
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
		return nil, errors.New("report category not found")
	}
	serviceTypes := r.Items.([]*domain.ServiceType)
	return serviceTypes[0], nil
}

func (s *ServiceTypeService) FindByNameAndServiceID(name string, serviceId int) (*domain.ServiceType, error) {
	r, err := s.Query(&models.ServiceTypesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceTypeFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: name,
			},
			ServiceID: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: serviceId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("report category not found")
	}
	serviceTypes := r.Items.([]*domain.ServiceType)
	return serviceTypes[0], nil
}
