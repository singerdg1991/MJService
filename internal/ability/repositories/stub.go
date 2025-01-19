package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/ability/domain"
	"github.com/hoitek/Maja-Service/internal/ability/models"
)

type AbilityRepositoryStub struct {
	Abilities []*domain.Ability
}

type abilityTestCondition struct {
	HasError bool
}

var UserTestCondition *abilityTestCondition = &abilityTestCondition{}

func NewAbilityRepositoryStub() *AbilityRepositoryStub {
	return &AbilityRepositoryStub{
		Abilities: []*domain.Ability{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *AbilityRepositoryStub) Query(dataModel *models.AbilitiesQueryRequestParams) ([]*domain.Ability, error) {
	var abilities []*domain.Ability
	for _, v := range r.Abilities {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			abilities = append(abilities, v)
			break
		}
	}
	return abilities, nil
}

func (r *AbilityRepositoryStub) Count(dataModel *models.AbilitiesQueryRequestParams) (int64, error) {
	var abilities []*domain.Ability
	for _, v := range r.Abilities {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			abilities = append(abilities, v)
			break
		}
	}
	return int64(len(abilities)), nil
}

func (r *AbilityRepositoryStub) Migrate() {
	// do stuff
}

func (r *AbilityRepositoryStub) Seed() {
	// do stuff
}

func (r *AbilityRepositoryStub) Create(payload *models.AbilitiesCreateRequestBody) (*domain.Ability, error) {
	panic("implement me")
}

func (r *AbilityRepositoryStub) Delete(payload *models.AbilitiesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *AbilityRepositoryStub) Update(payload *models.AbilitiesCreateRequestBody, name string) (*domain.Ability, error) {
	panic("implement me")
}

func (r *AbilityRepositoryStub) GetAbilitiesByIds(ids []int64) ([]*domain.Ability, error) {
	panic("implement me")
}
