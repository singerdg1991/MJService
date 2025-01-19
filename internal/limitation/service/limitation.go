package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/limitation/constants"
	"github.com/hoitek/Maja-Service/internal/limitation/domain"
	"github.com/hoitek/Maja-Service/internal/limitation/models"
	"github.com/hoitek/Maja-Service/internal/limitation/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type LimitationService struct {
	PostgresRepository ports.LimitationRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewLimitationService(pDB ports.LimitationRepositoryPostgresDB, m *storage.MinIO) LimitationService {
	go minio.SetupMinIOStorage(constants.LIMITATION_BUCKET_NAME, m)
	return LimitationService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *LimitationService) Query(q *models.LimitationsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying limitations", q)
	limitations, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.LimitationsQueryRequestParams{
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
	var items []*domain.Limitation
	for _, item := range limitations {
		items = append(items, &domain.Limitation{
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

func (s *LimitationService) Create(payload *models.LimitationsCreateRequestBody) (*domain.Limitation, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *LimitationService) Delete(payload *models.LimitationsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *LimitationService) Update(payload *models.LimitationsCreateRequestBody, id int64) (*domain.Limitation, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *LimitationService) GetLimitationsByIds(ids []int64) ([]*domain.Limitation, error) {
	return s.PostgresRepository.GetLimitationsByIds(ids)
}

func (s *LimitationService) FindByID(id int64) (*domain.Limitation, error) {
	r, err := s.Query(&models.LimitationsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("limitation not found")
	}
	limitations := r.Items.([]*domain.Limitation)
	return limitations[0], nil
}

func (s *LimitationService) FindByName(name string) (*domain.Limitation, error) {
	r, err := s.Query(&models.LimitationsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.LimitationFilterType{
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
		return nil, errors.New("limitation not found")
	}
	limitations := r.Items.([]*domain.Limitation)
	return limitations[0], nil
}
