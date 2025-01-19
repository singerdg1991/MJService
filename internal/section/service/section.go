package service

import (
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/section/constants"
	"github.com/hoitek/Maja-Service/internal/section/domain"
	"github.com/hoitek/Maja-Service/internal/section/models"
	"github.com/hoitek/Maja-Service/internal/section/ports"
	"github.com/hoitek/Maja-Service/storage"
	"math"
)

type SectionService struct {
	PostgresRepository ports.SectionRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewSectionService(pDB ports.SectionRepositoryPostgresDB, m *storage.MinIO) SectionService {
	go minio.SetupMinIOStorage(constants.SECTION_BUCKET_NAME, m)
	return SectionService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *SectionService) Query(q *models.SectionsQueryRequestParams) (*restypes.QueryResponse, error) {
	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 1, 1, q.Limit)

	sections, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(q)
	if err != nil {
		return nil, err
	}

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      sections,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *SectionService) Create(payload *models.SectionsCreateRequestBody) (*domain.Section, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *SectionService) Delete(payload *models.SectionsDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}

	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *SectionService) Update(payload *models.SectionsCreateRequestBody, id int) (*domain.Section, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *SectionService) GetSectionsByIds(ids []int64) ([]*domain.Section, error) {
	return s.PostgresRepository.GetSectionsByIds(ids)
}
