package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/servicegrade/constants"
	"github.com/hoitek/Maja-Service/internal/servicegrade/domain"
	"github.com/hoitek/Maja-Service/internal/servicegrade/models"
	"github.com/hoitek/Maja-Service/internal/servicegrade/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ServiceGradeService struct {
	PostgresRepository ports.ServiceGradeRepositoryPostgresDB
	MinioStorage       *storage.MinIO
}

func NewServiceGradeService(pDB ports.ServiceGradeRepositoryPostgresDB, m *storage.MinIO) ServiceGradeService {
	go minio.SetupMinIOStorage(constants.SERVICE_GRADE_BUCKET_NAME, m)
	return ServiceGradeService{
		PostgresRepository: pDB,
		MinioStorage:       m,
	}
}

func (s *ServiceGradeService) Query(q *models.ServiceGradesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying servicegrades", q)
	servicegrades, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ServiceGradesQueryRequestParams{
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
	var items []*domain.ServiceGrade
	for _, item := range servicegrades {
		items = append(items, &domain.ServiceGrade{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Grade:       item.Grade,
			Color:       item.Color,
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

func (s *ServiceGradeService) Create(payload *models.ServiceGradesCreateRequestBody) (*domain.ServiceGrade, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ServiceGradeService) Delete(payload *models.ServiceGradesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *ServiceGradeService) Update(payload *models.ServiceGradesCreateRequestBody, id int64) (*domain.ServiceGrade, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ServiceGradeService) GetServiceGradesByIds(ids []int64) ([]*domain.ServiceGrade, error) {
	return s.PostgresRepository.GetServiceGradesByIds(ids)
}

func (s *ServiceGradeService) FindByID(id int64) (*domain.ServiceGrade, error) {
	r, err := s.Query(&models.ServiceGradesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	servicegrades := r.Items.([]*domain.ServiceGrade)
	return servicegrades[0], nil
}

func (s *ServiceGradeService) FindByName(name string) (*domain.ServiceGrade, error) {
	r, err := s.Query(&models.ServiceGradesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ServiceGradeFilterType{
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
		return nil, errors.New("language skill not found")
	}
	servicegrades := r.Items.([]*domain.ServiceGrade)
	return servicegrades[0], nil
}
