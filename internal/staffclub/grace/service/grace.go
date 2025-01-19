package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/constants"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type GraceService struct {
	PostgresRepository ports.GraceRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewGraceService(pDB ports.GraceRepositoryPostgresDB, m *storage.MinIO) GraceService {
	go minio.SetupMinIOStorage(constants.GRACE_BUCKET_NAME, m)
	return GraceService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *GraceService) Query(q *models.GracesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying graces", q)
	graces, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.GracesQueryRequestParams{
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
		Items:      graces,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *GraceService) Create(payload *models.GracesCreateRequestBody) (*domain.Grace, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *GraceService) Delete(payload *models.GracesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *GraceService) Update(payload *models.GracesCreateRequestBody, id int64) (*domain.Grace, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *GraceService) GetGracesByIds(ids []int64) ([]*domain.Grace, error) {
	return s.PostgresRepository.GetGracesByIds(ids)
}

func (s *GraceService) FindByID(id int64) (*domain.Grace, error) {
	r, err := s.Query(&models.GracesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("grace not found")
	}
	graces := r.Items.([]*domain.Grace)
	return graces[0], nil
}

func (s *GraceService) FindByGraceNumber(graceNumber int64) (*domain.Grace, error) {
	r, err := s.Query(&models.GracesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.GraceFilterType{
			GraceNumber: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: int(graceNumber),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("grace not found")
	}
	graces, ok := r.Items.([]*domain.Grace)
	if !ok {
		return nil, errors.New("grace not found")
	}
	return graces[0], nil
}
