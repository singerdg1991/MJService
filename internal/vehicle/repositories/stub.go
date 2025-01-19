package repositories

import (
	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
	"github.com/hoitek/Maja-Service/internal/vehicle/models"
)

type VehicleRepositoryStub struct {
	Vehicles []*domain.Vehicle
}

type vehicleTestCondition struct {
	HasError bool
}

var UserTestCondition *vehicleTestCondition = &vehicleTestCondition{}

func NewVehicleRepositoryStub() *VehicleRepositoryStub {
	return &VehicleRepositoryStub{
		Vehicles: []*domain.Vehicle{
			{
				ID: 1,
			},
		},
	}
}

func (r *VehicleRepositoryStub) Query(dataModel *models.VehiclesQueryRequestParams) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	for _, v := range r.Vehicles {
		if v.ID == uint(dataModel.ID) {
			vehicles = append(vehicles, v)
			break
		}
	}
	return vehicles, nil
}

func (r *VehicleRepositoryStub) Count(dataModel *models.VehiclesQueryRequestParams) (int64, error) {
	var vehicles []*domain.Vehicle
	for _, v := range r.Vehicles {
		if v.ID == uint(dataModel.ID) {
			vehicles = append(vehicles, v)
			break
		}
	}
	return int64(len(vehicles)), nil
}

func (r *VehicleRepositoryStub) Migrate() {
	// do stuff
}

func (r *VehicleRepositoryStub) Seed() {
	// do stuff
}

func (r *VehicleRepositoryStub) Create(payload *models.VehiclesCreateRequestBody) (*domain.Vehicle, error) {
	panic("implement me")
}

func (r *VehicleRepositoryStub) Delete(payload *models.VehiclesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *VehicleRepositoryStub) Update(payload *models.VehiclesCreateRequestBody, id int) (*domain.Vehicle, error) {
	panic("implement me")
}
