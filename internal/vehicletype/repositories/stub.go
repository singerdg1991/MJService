package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/vehicletype/domain"
	"github.com/hoitek/Maja-Service/internal/vehicletype/models"
)

type VehicleTypeRepositoryStub struct {
	VehicleTypes []*domain.VehicleType
}

type vehicletypeTestCondition struct {
	HasError bool
}

var UserTestCondition *vehicletypeTestCondition = &vehicletypeTestCondition{}

func NewVehicleTypeRepositoryStub() *VehicleTypeRepositoryStub {
	return &VehicleTypeRepositoryStub{
		VehicleTypes: []*domain.VehicleType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *VehicleTypeRepositoryStub) Query(dataModel *models.VehicleTypesQueryRequestParams) ([]*domain.VehicleType, error) {
	var vehicletypes []*domain.VehicleType
	for _, v := range r.VehicleTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			vehicletypes = append(vehicletypes, v)
			break
		}
	}
	return vehicletypes, nil
}

func (r *VehicleTypeRepositoryStub) Count(dataModel *models.VehicleTypesQueryRequestParams) (int64, error) {
	var vehicletypes []*domain.VehicleType
	for _, v := range r.VehicleTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			vehicletypes = append(vehicletypes, v)
			break
		}
	}
	return int64(len(vehicletypes)), nil
}

func (r *VehicleTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *VehicleTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *VehicleTypeRepositoryStub) Create(payload *models.VehicleTypesCreateRequestBody) (*domain.VehicleType, error) {
	panic("implement me")
}

func (r *VehicleTypeRepositoryStub) Delete(payload *models.VehicleTypesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *VehicleTypeRepositoryStub) Update(payload *models.VehicleTypesCreateRequestBody, name string) (*domain.VehicleType, error) {
	panic("implement me")
}
