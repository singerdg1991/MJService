package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/limitation/domain"
	"github.com/hoitek/Maja-Service/internal/limitation/models"
)

type LimitationRepositoryStub struct {
	Limitations []*domain.Limitation
}

type limitationTestCondition struct {
	HasError bool
}

var UserTestCondition *limitationTestCondition = &limitationTestCondition{}

func NewLimitationRepositoryStub() *LimitationRepositoryStub {
	return &LimitationRepositoryStub{
		Limitations: []*domain.Limitation{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *LimitationRepositoryStub) Query(dataModel *models.LimitationsQueryRequestParams) ([]*domain.Limitation, error) {
	var limitations []*domain.Limitation
	for _, v := range r.Limitations {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			limitations = append(limitations, v)
			break
		}
	}
	return limitations, nil
}

func (r *LimitationRepositoryStub) Count(dataModel *models.LimitationsQueryRequestParams) (int64, error) {
	var limitations []*domain.Limitation
	for _, v := range r.Limitations {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			limitations = append(limitations, v)
			break
		}
	}
	return int64(len(limitations)), nil
}

func (r *LimitationRepositoryStub) Migrate() {
	// do stuff
}

func (r *LimitationRepositoryStub) Seed() {
	// do stuff
}

func (r *LimitationRepositoryStub) Create(payload *models.LimitationsCreateRequestBody) (*domain.Limitation, error) {
	panic("implement me")
}

func (r *LimitationRepositoryStub) Delete(payload *models.LimitationsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *LimitationRepositoryStub) Update(payload *models.LimitationsCreateRequestBody, id int64) (*domain.Limitation, error) {
	panic("implement me")
}

func (r *LimitationRepositoryStub) GetLimitationsByIds(ids []int64) ([]*domain.Limitation, error) {
	panic("implement me")
}
