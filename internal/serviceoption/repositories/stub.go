package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/serviceoption/domain"
	"github.com/hoitek/Maja-Service/internal/serviceoption/models"
)

type ServiceOptionRepositoryStub struct {
	ServiceOptions []*domain.ServiceOption
}

type serviceTypeTestCondition struct {
	HasError bool
}

var UserTestCondition *serviceTypeTestCondition = &serviceTypeTestCondition{}

func NewServiceOptionRepositoryStub() *ServiceOptionRepositoryStub {
	return &ServiceOptionRepositoryStub{
		ServiceOptions: []*domain.ServiceOption{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ServiceOptionRepositoryStub) Query(dataModel *models.ServiceOptionsQueryRequestParams) ([]*domain.ServiceOption, error) {
	var reports []*domain.ServiceOption
	for _, v := range r.ServiceOptions {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			reports = append(reports, v)
			break
		}
	}
	return reports, nil
}

func (r *ServiceOptionRepositoryStub) Count(dataModel *models.ServiceOptionsQueryRequestParams) (int64, error) {
	var reports []*domain.ServiceOption
	for _, v := range r.ServiceOptions {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			reports = append(reports, v)
			break
		}
	}
	return int64(len(reports)), nil
}

func (r *ServiceOptionRepositoryStub) Migrate() {
	// do stuff
}

func (r *ServiceOptionRepositoryStub) Seed() {
	// do stuff
}

func (r *ServiceOptionRepositoryStub) Create(payload *models.ServiceOptionsCreateRequestBody) (*domain.ServiceOption, error) {
	panic("implement me")
}

func (r *ServiceOptionRepositoryStub) Delete(payload *models.ServiceOptionsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ServiceOptionRepositoryStub) Update(payload *models.ServiceOptionsCreateRequestBody, id int64) (*domain.ServiceOption, error) {
	panic("implement me")
}

func (r *ServiceOptionRepositoryStub) GetServiceOptionsByIds(ids []int64) ([]*domain.ServiceOption, error) {
	panic("implement me")
}
