package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/ability/constants"
	"github.com/hoitek/Maja-Service/internal/ability/domain"
	"github.com/hoitek/Maja-Service/internal/ability/models"
	"github.com/hoitek/Maja-Service/internal/ability/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type AbilityService struct {
	PostgresRepository ports.AbilityRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewAbilityService(pDB ports.AbilityRepositoryPostgresDB, m *storage.MinIO) AbilityService {
	go minio.SetupMinIOStorage(constants.ABILITY_BUCKET_NAME, m)
	return AbilityService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *AbilityService) Query(q *models.AbilitiesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying abilities", q)
	abilities, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.AbilitiesQueryRequestParams{
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
	var items []*domain.Ability
	for _, item := range abilities {
		items = append(items, &domain.Ability{
			ID:   item.ID,
			Name: item.Name,
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

func (s *AbilityService) Create(payload *models.AbilitiesCreateRequestBody) (*domain.Ability, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *AbilityService) Delete(payload *models.AbilitiesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *AbilityService) Update(payload *models.AbilitiesCreateRequestBody, name string) (*domain.Ability, error) {
	return s.PostgresRepository.Update(payload, name)
}

func (s *AbilityService) GetAbilitiesByIds(ids []int64) ([]*domain.Ability, error) {
	return s.PostgresRepository.GetAbilitiesByIds(ids)
}
