package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/models"
)

type GraceRepositoryStub struct {
	Graces []*domain.Grace
}

type graceTestCondition struct {
	HasError bool
}

var UserTestCondition *graceTestCondition = &graceTestCondition{}

func NewGraceRepositoryStub() *GraceRepositoryStub {
	return &GraceRepositoryStub{
		Graces: []*domain.Grace{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *GraceRepositoryStub) Query(dataModel *models.GracesQueryRequestParams) ([]*domain.Grace, error) {
	var graces []*domain.Grace
	for _, v := range r.Graces {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			graces = append(graces, v)
			break
		}
	}
	return graces, nil
}

func (r *GraceRepositoryStub) Count(dataModel *models.GracesQueryRequestParams) (int64, error) {
	var graces []*domain.Grace
	for _, v := range r.Graces {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			graces = append(graces, v)
			break
		}
	}
	return int64(len(graces)), nil
}

func (r *GraceRepositoryStub) Migrate() {
	// do stuff
}

func (r *GraceRepositoryStub) Seed() {
	// do stuff
}

func (r *GraceRepositoryStub) Create(payload *models.GracesCreateRequestBody) (*domain.Grace, error) {
	panic("implement me")
}

func (r *GraceRepositoryStub) Delete(payload *models.GracesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *GraceRepositoryStub) Update(payload *models.GracesCreateRequestBody, id int64) (*domain.Grace, error) {
	panic("implement me")
}

func (r *GraceRepositoryStub) GetGracesByIds(ids []int64) ([]*domain.Grace, error) {
	panic("implement me")
}
