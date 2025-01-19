package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/models"
)

type WarningRepositoryStub struct {
	Warnings []*domain.Warning
}

type warningTestCondition struct {
	HasError bool
}

var UserTestCondition *warningTestCondition = &warningTestCondition{}

func NewWarningRepositoryStub() *WarningRepositoryStub {
	return &WarningRepositoryStub{
		Warnings: []*domain.Warning{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *WarningRepositoryStub) Query(dataModel *models.WarningsQueryRequestParams) ([]*domain.Warning, error) {
	var warnings []*domain.Warning
	for _, v := range r.Warnings {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			warnings = append(warnings, v)
			break
		}
	}
	return warnings, nil
}

func (r *WarningRepositoryStub) Count(dataModel *models.WarningsQueryRequestParams) (int64, error) {
	var warnings []*domain.Warning
	for _, v := range r.Warnings {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			warnings = append(warnings, v)
			break
		}
	}
	return int64(len(warnings)), nil
}

func (r *WarningRepositoryStub) Migrate() {
	// do stuff
}

func (r *WarningRepositoryStub) Seed() {
	// do stuff
}

func (r *WarningRepositoryStub) Create(payload *models.WarningsCreateRequestBody) (*domain.Warning, error) {
	panic("implement me")
}

func (r *WarningRepositoryStub) Delete(payload *models.WarningsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *WarningRepositoryStub) Update(payload *models.WarningsCreateRequestBody, id int64) (*domain.Warning, error) {
	panic("implement me")
}

func (r *WarningRepositoryStub) GetWarningsByIds(ids []int64) ([]*domain.Warning, error) {
	panic("implement me")
}
