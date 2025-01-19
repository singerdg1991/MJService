package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/constants"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type AttentionService struct {
	PostgresRepository ports.AttentionRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewAttentionService(pDB ports.AttentionRepositoryPostgresDB, m *storage.MinIO) AttentionService {
	go minio.SetupMinIOStorage(constants.ATTENTION_BUCKET_NAME, m)
	return AttentionService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *AttentionService) Query(q *models.AttentionsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying attentions", q)
	attentions, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.AttentionsQueryRequestParams{
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
		Items:      attentions,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *AttentionService) Create(payload *models.AttentionsCreateRequestBody) (*domain.Attention, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *AttentionService) Delete(payload *models.AttentionsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *AttentionService) Update(payload *models.AttentionsCreateRequestBody, id int64) (*domain.Attention, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *AttentionService) GetAttentionsByIds(ids []int64) ([]*domain.Attention, error) {
	return s.PostgresRepository.GetAttentionsByIds(ids)
}

func (s *AttentionService) FindByID(id int64) (*domain.Attention, error) {
	r, err := s.Query(&models.AttentionsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("attention not found")
	}
	attentions := r.Items.([]*domain.Attention)
	return attentions[0], nil
}

func (s *AttentionService) FindByAttentionNumber(attentionNumber int64) (*domain.Attention, error) {
	r, err := s.Query(&models.AttentionsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.AttentionFilterType{
			AttentionNumber: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: int(attentionNumber),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("attention not found")
	}
	attentions, ok := r.Items.([]*domain.Attention)
	if !ok {
		return nil, errors.New("attention not found")
	}
	return attentions[0], nil
}
