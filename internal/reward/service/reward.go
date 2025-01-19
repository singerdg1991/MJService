package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/reward/constants"
	"github.com/hoitek/Maja-Service/internal/reward/domain"
	"github.com/hoitek/Maja-Service/internal/reward/models"
	"github.com/hoitek/Maja-Service/internal/reward/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type RewardService struct {
	PostgresRepository ports.RewardRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewRewardService(pDB ports.RewardRepositoryPostgresDB, m *storage.MinIO) RewardService {
	go minio.SetupMinIOStorage(constants.REWARD_BUCKET_NAME, m)
	return RewardService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *RewardService) Query(q *models.RewardsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying rewards", q)
	rewards, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.RewardsQueryRequestParams{
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
	var items []*domain.Reward
	for _, item := range rewards {
		items = append(items, &domain.Reward{
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

func (s *RewardService) Create(payload *models.RewardsCreateRequestBody) (*domain.Reward, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *RewardService) Delete(payload *models.RewardsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *RewardService) Update(payload *models.RewardsCreateRequestBody, id int64) (*domain.Reward, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *RewardService) GetRewardsByIds(ids []int64) ([]*domain.Reward, error) {
	return s.PostgresRepository.GetRewardsByIds(ids)
}

func (s *RewardService) FindByID(id int64) (*domain.Reward, error) {
	r, err := s.Query(&models.RewardsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("reward not found")
	}
	rewards, ok := r.Items.([]*domain.Reward)
	if !ok {
		return nil, errors.New("reward not found")
	}
	return rewards[0], nil
}

func (s *RewardService) FindByName(name string) (*domain.Reward, error) {
	r, err := s.Query(&models.RewardsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.RewardFilterType{
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
		return nil, errors.New("reward not found")
	}
	rewards := r.Items.([]*domain.Reward)
	return rewards[0], nil
}
