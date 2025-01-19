package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/geartype/domain"
	"github.com/hoitek/Maja-Service/internal/geartype/models"
)

type GearTypeRepositoryStub struct {
	GearTypes []*domain.GearType
}

type geartypeTestCondition struct {
	HasError bool
}

var UserTestCondition *geartypeTestCondition = &geartypeTestCondition{}

func NewGearTypeRepositoryStub() *GearTypeRepositoryStub {
	return &GearTypeRepositoryStub{
		GearTypes: []*domain.GearType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *GearTypeRepositoryStub) Query(dataModel *models.GearTypesQueryRequestParams) ([]*domain.GearType, error) {
	var geartypes []*domain.GearType
	for _, v := range r.GearTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			geartypes = append(geartypes, v)
			break
		}
	}
	return geartypes, nil
}

func (r *GearTypeRepositoryStub) Count(dataModel *models.GearTypesQueryRequestParams) (int64, error) {
	var geartypes []*domain.GearType
	for _, v := range r.GearTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			geartypes = append(geartypes, v)
			break
		}
	}
	return int64(len(geartypes)), nil
}

func (r *GearTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *GearTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *GearTypeRepositoryStub) Create(payload *models.GearTypesCreateRequestBody) (*domain.GearType, error) {
	panic("implement me")
}

func (r *GearTypeRepositoryStub) Delete(payload *models.GearTypesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *GearTypeRepositoryStub) Update(payload *models.GearTypesCreateRequestBody, name string) (*domain.GearType, error) {
	panic("implement me")
}
