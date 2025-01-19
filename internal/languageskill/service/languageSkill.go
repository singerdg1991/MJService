package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/languageskill/constants"
	"github.com/hoitek/Maja-Service/internal/languageskill/domain"
	"github.com/hoitek/Maja-Service/internal/languageskill/models"
	"github.com/hoitek/Maja-Service/internal/languageskill/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type LanguageSkillService struct {
	PostgresRepository ports.LanguageSkillRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewLanguageSkillService(pDB ports.LanguageSkillRepositoryPostgresDB, m *storage.MinIO) LanguageSkillService {
	go minio.SetupMinIOStorage(constants.LANGUAGE_SKILL_BUCKET_NAME, m)
	return LanguageSkillService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *LanguageSkillService) Query(q *models.LanguageSkillsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying languageskills", q)
	languageSkills, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.LanguageSkillsQueryRequestParams{
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
	var items []*domain.LanguageSkill
	for _, item := range languageSkills {
		items = append(items, &domain.LanguageSkill{
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

func (s *LanguageSkillService) Create(payload *models.LanguageSkillsCreateRequestBody) (*domain.LanguageSkill, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *LanguageSkillService) Delete(payload *models.LanguageSkillsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *LanguageSkillService) Update(payload *models.LanguageSkillsCreateRequestBody, id int64) (*domain.LanguageSkill, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *LanguageSkillService) GetLanguageSkillsByIds(ids []int64) ([]*domain.LanguageSkill, error) {
	return s.PostgresRepository.GetLanguageSkillsByIds(ids)
}

func (s *LanguageSkillService) FindByID(id int64) (*domain.LanguageSkill, error) {
	r, err := s.Query(&models.LanguageSkillsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	languageSkills := r.Items.([]*domain.LanguageSkill)
	return languageSkills[0], nil
}

func (s *LanguageSkillService) FindByName(name string) (*domain.LanguageSkill, error) {
	r, err := s.Query(&models.LanguageSkillsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.LanguageSkillFilterType{
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
	languageSkills := r.Items.([]*domain.LanguageSkill)
	return languageSkills[0], nil
}
