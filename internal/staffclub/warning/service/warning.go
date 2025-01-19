package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/constants"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type WarningService struct {
	PostgresRepository ports.WarningRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewWarningService(pDB ports.WarningRepositoryPostgresDB, m *storage.MinIO) WarningService {
	go minio.SetupMinIOStorage(constants.WARNING_BUCKET_NAME, m)
	return WarningService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *WarningService) Query(q *models.WarningsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying warnings", q)
	warnings, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.WarningsQueryRequestParams{
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

	return &restypes.QueryResponse{
		Items:      warnings,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *WarningService) Create(payload *models.WarningsCreateRequestBody) (*domain.Warning, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *WarningService) Delete(payload *models.WarningsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *WarningService) Update(payload *models.WarningsCreateRequestBody, id int64) (*domain.Warning, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *WarningService) GetWarningsByIds(ids []int64) ([]*domain.Warning, error) {
	return s.PostgresRepository.GetWarningsByIds(ids)
}

func (s *WarningService) FindByID(id int64) (*domain.Warning, error) {
	r, err := s.Query(&models.WarningsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("warning not found")
	}
	warnings := r.Items.([]*domain.Warning)
	return warnings[0], nil
}

func (s *WarningService) FindByWarningNumber(warningNumber int64) (*domain.Warning, error) {
	r, err := s.Query(&models.WarningsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.WarningFilterType{
			WarningNumber: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: int(warningNumber),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("warning not found")
	}
	warnings, ok := r.Items.([]*domain.Warning)
	if !ok {
		return nil, errors.New("warning not found")
	}
	return warnings[0], nil
}
