package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/serviceoption/constants"
	"github.com/hoitek/Maja-Service/internal/serviceoption/domain"
	"github.com/hoitek/Maja-Service/internal/serviceoption/models"
	"github.com/hoitek/Maja-Service/internal/serviceoption/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ServiceService struct {
	PostgresRepository ports.ServiceOptionRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewServiceService(pDB ports.ServiceOptionRepositoryPostgresDB, m *storage.MinIO) ServiceService {
	go minio.SetupMinIOStorage(constants.SERVICE_OPTION_BUCKET_NAME, m)
	return ServiceService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ServiceService) Query(q *models.ServiceOptionsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying reports", q)
	reports, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ServiceOptionsQueryRequestParams{
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
	var items []*domain.ServiceOption
	for _, item := range reports {
		items = append(items, &domain.ServiceOption{
			ID:            item.ID,
			ServiceTypeID: item.ServiceTypeID,
			ServiceType:   item.ServiceType,
			Name:          item.Name,
			Description:   item.Description,
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

func (s *ServiceService) Create(payload *models.ServiceOptionsCreateRequestBody) (*domain.ServiceOption, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ServiceService) Delete(payload *models.ServiceOptionsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *ServiceService) Update(payload *models.ServiceOptionsCreateRequestBody, id int64) (*domain.ServiceOption, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ServiceService) GetServiceOptionsByIds(ids []int64) ([]*domain.ServiceOption, error) {
	return s.PostgresRepository.GetServiceOptionsByIds(ids)
}

func (s *ServiceService) FindByID(id int64) (*domain.ServiceOption, error) {
	r, err := s.Query(&models.ServiceOptionsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("serviceoption not found")
	}
	reports := r.Items.([]*domain.ServiceOption)
	return reports[0], nil
}

func (s *ServiceService) FindByName(name string) (*domain.ServiceOption, error) {
	r, err := s.Query(&models.ServiceOptionsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceOptionFilterType{
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
		return nil, errors.New("serviceoption not found")
	}
	reports := r.Items.([]*domain.ServiceOption)
	return reports[0], nil
}

func (s *ServiceService) FindByNameAndServiceTypeID(name string, serviceTypeId int) (*domain.ServiceOption, error) {
	r, err := s.Query(&models.ServiceOptionsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceOptionFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: name,
			},
			ServiceTypeID: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: serviceTypeId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("serviceoption not found")
	}
	reports := r.Items.([]*domain.ServiceOption)
	return reports[0], nil
}
