package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/medicine/domain"
	"github.com/hoitek/Maja-Service/internal/medicine/models"
)

type MedicineRepositoryStub struct {
	Medicines []*domain.Medicine
}

type medicineTestCondition struct {
	HasError bool
}

var UserTestCondition *medicineTestCondition = &medicineTestCondition{}

func NewMedicineRepositoryStub() *MedicineRepositoryStub {
	return &MedicineRepositoryStub{
		Medicines: []*domain.Medicine{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *MedicineRepositoryStub) Query(dataModel *models.MedicinesQueryRequestParams) ([]*domain.Medicine, error) {
	var medicines []*domain.Medicine
	for _, v := range r.Medicines {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			medicines = append(medicines, v)
			break
		}
	}
	return medicines, nil
}

func (r *MedicineRepositoryStub) Count(dataModel *models.MedicinesQueryRequestParams) (int64, error) {
	var medicines []*domain.Medicine
	for _, v := range r.Medicines {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			medicines = append(medicines, v)
			break
		}
	}
	return int64(len(medicines)), nil
}

func (r *MedicineRepositoryStub) Migrate() {
	// do stuff
}

func (r *MedicineRepositoryStub) Seed() {
	// do stuff
}

func (r *MedicineRepositoryStub) Create(payload *models.MedicinesCreateRequestBody) (*domain.Medicine, error) {
	panic("implement me")
}

func (r *MedicineRepositoryStub) Delete(payload *models.MedicinesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *MedicineRepositoryStub) Update(payload *models.MedicinesCreateRequestBody, id int64) (*domain.Medicine, error) {
	panic("implement me")
}

func (r *MedicineRepositoryStub) GetMedicinesByIds(ids []int64) ([]*domain.Medicine, error) {
	panic("implement me")
}
