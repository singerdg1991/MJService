package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/reward/domain"
	"github.com/hoitek/Maja-Service/internal/reward/models"
)

type RewardRepositoryStub struct {
	Rewards []*domain.Reward
}

type rewardTestCondition struct {
	HasError bool
}

var UserTestCondition *rewardTestCondition = &rewardTestCondition{}

func NewRewardRepositoryStub() *RewardRepositoryStub {
	return &RewardRepositoryStub{
		Rewards: []*domain.Reward{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *RewardRepositoryStub) Query(dataModel *models.RewardsQueryRequestParams) ([]*domain.Reward, error) {
	var rewards []*domain.Reward
	for _, v := range r.Rewards {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			rewards = append(rewards, v)
			break
		}
	}
	return rewards, nil
}

func (r *RewardRepositoryStub) Count(dataModel *models.RewardsQueryRequestParams) (int64, error) {
	var rewards []*domain.Reward
	for _, v := range r.Rewards {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			rewards = append(rewards, v)
			break
		}
	}
	return int64(len(rewards)), nil
}

func (r *RewardRepositoryStub) Migrate() {
	// do stuff
}

func (r *RewardRepositoryStub) Seed() {
	// do stuff
}

func (r *RewardRepositoryStub) Create(payload *models.RewardsCreateRequestBody) (*domain.Reward, error) {
	panic("implement me")
}

func (r *RewardRepositoryStub) Delete(payload *models.RewardsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *RewardRepositoryStub) Update(payload *models.RewardsCreateRequestBody, id int64) (*domain.Reward, error) {
	panic("implement me")
}

func (r *RewardRepositoryStub) GetRewardsByIds(ids []int64) ([]*domain.Reward, error) {
	panic("implement me")
}
