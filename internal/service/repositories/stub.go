package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/service/domain"
	"github.com/hoitek/Maja-Service/internal/service/models"
)

type ServiceRepositoryStub struct {
	Services []*domain.Service
}

type serviceTestCondition struct {
	HasError bool
}

var UserTestCondition *serviceTestCondition = &serviceTestCondition{}

func NewServiceRepositoryStub() *ServiceRepositoryStub {
	return &ServiceRepositoryStub{
		Services: []*domain.Service{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ServiceRepositoryStub) Query(dataModel *models.ServicesQueryRequestParams) ([]*domain.Service, error) {
	var services []*domain.Service
	for _, v := range r.Services {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			services = append(services, v)
			break
		}
	}
	return services, nil
}

func (r *ServiceRepositoryStub) Count(dataModel *models.ServicesQueryRequestParams) (int64, error) {
	var services []*domain.Service
	for _, v := range r.Services {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			services = append(services, v)
			break
		}
	}
	return int64(len(services)), nil
}

func (r *ServiceRepositoryStub) Migrate() {
	// do stuff
}

func (r *ServiceRepositoryStub) Seed() {
	// do stuff
}

func (r *ServiceRepositoryStub) Create(payload *models.ServicesCreateRequestBody) (*domain.Service, error) {
	panic("implement me")
}

func (r *ServiceRepositoryStub) Delete(payload *models.ServicesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ServiceRepositoryStub) Update(payload *models.ServicesCreateRequestBody, id int64) (*domain.Service, error) {
	panic("implement me")
}

func (r *ServiceRepositoryStub) GetServicesByIds(ids []int64) ([]*domain.Service, error) {
	panic("implement me")
}
