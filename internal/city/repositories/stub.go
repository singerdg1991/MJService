package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/city/domain"
	"github.com/hoitek/Maja-Service/internal/city/models"
)

type CityRepositoryStub struct {
	Cities []*domain.City
}

type cityTestCondition struct {
	HasError bool
}

var UserTestCondition *cityTestCondition = &cityTestCondition{}

func NewCityRepositoryStub() *CityRepositoryStub {
	return &CityRepositoryStub{
		Cities: []*domain.City{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *CityRepositoryStub) Query(dataModel *models.CitiesQueryRequestParams) ([]*domain.City, error) {
	var cities []*domain.City
	for _, v := range r.Cities {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			cities = append(cities, v)
			break
		}
	}
	return cities, nil
}

func (r *CityRepositoryStub) Count(dataModel *models.CitiesQueryRequestParams) (int64, error) {
	var cities []*domain.City
	for _, v := range r.Cities {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			cities = append(cities, v)
			break
		}
	}
	return int64(len(cities)), nil
}

func (r *CityRepositoryStub) Migrate() {
	// do stuff
}

func (r *CityRepositoryStub) Seed() {
	// do stuff
}

func (r *CityRepositoryStub) Create(payload *models.CitiesCreateRequestBody) (*domain.City, error) {
	panic("implement me")
}

func (r *CityRepositoryStub) Delete(payload *models.CitiesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CityRepositoryStub) Update(payload *models.CitiesCreateRequestBody, name string) (*domain.City, error) {
	panic("implement me")
}
