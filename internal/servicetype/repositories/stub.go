package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/servicetype/domain"
	"github.com/hoitek/Maja-Service/internal/servicetype/models"
)

type ServiceTypeRepositoryStub struct {
	ServiceTypes []*domain.ServiceType
}

type serviceTypeTestCondition struct {
	HasError bool
}

var UserTestCondition *serviceTypeTestCondition = &serviceTypeTestCondition{}

func NewServiceTypeRepositoryStub() *ServiceTypeRepositoryStub {
	return &ServiceTypeRepositoryStub{
		ServiceTypes: []*domain.ServiceType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ServiceTypeRepositoryStub) Query(dataModel *models.ServiceTypesQueryRequestParams) ([]*domain.ServiceType, error) {
	var serviceTypes []*domain.ServiceType
	for _, v := range r.ServiceTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			serviceTypes = append(serviceTypes, v)
			break
		}
	}
	return serviceTypes, nil
}

func (r *ServiceTypeRepositoryStub) Count(dataModel *models.ServiceTypesQueryRequestParams) (int64, error) {
	var serviceTypes []*domain.ServiceType
	for _, v := range r.ServiceTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			serviceTypes = append(serviceTypes, v)
			break
		}
	}
	return int64(len(serviceTypes)), nil
}

func (r *ServiceTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *ServiceTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *ServiceTypeRepositoryStub) Create(payload *models.ServiceTypesCreateRequestBody) (*domain.ServiceType, error) {
	panic("implement me")
}

func (r *ServiceTypeRepositoryStub) Delete(payload *models.ServiceTypesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ServiceTypeRepositoryStub) Update(payload *models.ServiceTypesCreateRequestBody, id int64) (*domain.ServiceType, error) {
	panic("implement me")
}

func (r *ServiceTypeRepositoryStub) GetServiceTypesByIds(ids []int64) ([]*domain.ServiceType, error) {
	panic("implement me")
}
