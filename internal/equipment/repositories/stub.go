package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/equipment/domain"
	"github.com/hoitek/Maja-Service/internal/equipment/models"
)

type EquipmentRepositoryStub struct {
	Equipments []*domain.Equipment
}

type equipmentTestCondition struct {
	HasError bool
}

var UserTestCondition *equipmentTestCondition = &equipmentTestCondition{}

func NewEquipmentRepositoryStub() *EquipmentRepositoryStub {
	return &EquipmentRepositoryStub{
		Equipments: []*domain.Equipment{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *EquipmentRepositoryStub) Query(dataModel *models.EquipmentsQueryRequestParams) ([]*domain.Equipment, error) {
	var equipments []*domain.Equipment
	for _, v := range r.Equipments {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			equipments = append(equipments, v)
			break
		}
	}
	return equipments, nil
}

func (r *EquipmentRepositoryStub) Count(dataModel *models.EquipmentsQueryRequestParams) (int64, error) {
	var equipments []*domain.Equipment
	for _, v := range r.Equipments {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			equipments = append(equipments, v)
			break
		}
	}
	return int64(len(equipments)), nil
}

func (r *EquipmentRepositoryStub) Migrate() {
	// do stuff
}

func (r *EquipmentRepositoryStub) Seed() {
	// do stuff
}

func (r *EquipmentRepositoryStub) Create(payload *models.EquipmentsCreateRequestBody) (*domain.Equipment, error) {
	panic("implement me")
}

func (r *EquipmentRepositoryStub) Delete(payload *models.EquipmentsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *EquipmentRepositoryStub) Update(payload *models.EquipmentsCreateRequestBody, id int64) (*domain.Equipment, error) {
	panic("implement me")
}

func (r *EquipmentRepositoryStub) GetEquipmentsByIds(ids []int64) ([]*domain.Equipment, error) {
	panic("implement me")
}
