package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/reward/domain"
	"github.com/hoitek/Maja-Service/internal/reward/models"
)

type RewardService interface {
	Query(dataModel *models.RewardsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.RewardsCreateRequestBody) (*domain.Reward, error)
	Delete(payload *models.RewardsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.RewardsCreateRequestBody, id int64) (*domain.Reward, error)
	GetRewardsByIds(ids []int64) ([]*domain.Reward, error)
	FindByID(id int64) (*domain.Reward, error)
	FindByName(name string) (*domain.Reward, error)
}

type RewardRepositoryPostgresDB interface {
	Query(dataModel *models.RewardsQueryRequestParams) ([]*domain.Reward, error)
	Count(dataModel *models.RewardsQueryRequestParams) (int64, error)
	Create(payload *models.RewardsCreateRequestBody) (*domain.Reward, error)
	Delete(payload *models.RewardsDeleteRequestBody) ([]int64, error)
	Update(payload *models.RewardsCreateRequestBody, id int64) (*domain.Reward, error)
	GetRewardsByIds(ids []int64) ([]*domain.Reward, error)
}
