package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/evaluation/constants"
	"github.com/hoitek/Maja-Service/internal/evaluation/domain"
	"github.com/hoitek/Maja-Service/internal/evaluation/models"
	"github.com/hoitek/Maja-Service/internal/evaluation/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type EvaluationService struct {
	PostgresRepository ports.EvaluationRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewEvaluationService(pDB ports.EvaluationRepositoryPostgresDB, m *storage.MinIO) EvaluationService {
	go minio.SetupMinIOStorage(constants.EVALUATION_BUCKET_NAME, m)
	return EvaluationService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *EvaluationService) Query(q *models.EvaluationsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying evaluations", q)
	evaluations, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.EvaluationsQueryRequestParams{
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
		Items:      evaluations,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *EvaluationService) Create(payload *models.EvaluationsCreateRequestBody) (*domain.Evaluation, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *EvaluationService) Delete(payload *models.EvaluationsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *EvaluationService) Update(payload *models.EvaluationsCreateRequestBody, id int64) (*domain.Evaluation, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *EvaluationService) GetEvaluationsByIds(ids []int64) ([]*domain.Evaluation, error) {
	return s.PostgresRepository.GetEvaluationsByIds(ids)
}

func (s *EvaluationService) FindByID(id int64) (*domain.Evaluation, error) {
	r, err := s.Query(&models.EvaluationsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("evaluation not found")
	}
	evaluations, ok := r.Items.([]*domain.Evaluation)
	if !ok {
		return nil, errors.New("evaluation not found")
	}
	return evaluations[0], nil
}

func (s *EvaluationService) FindByTitle(title string) (*domain.Evaluation, error) {
	r, err := s.Query(&models.EvaluationsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.EvaluationFilterType{
			Title: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: title,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("evaluation not found")
	}
	evaluations := r.Items.([]*domain.Evaluation)
	return evaluations[0], nil
}
