package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/shifttype/domain"
	"github.com/hoitek/Maja-Service/internal/shifttype/models"
)

type ShiftTypeRepositoryStub struct {
	ShiftTypes []*domain.ShiftType
}

type ShiftTypeTestCondition struct {
	HasError bool
}

func NewShiftTypeRepositoryStub() *ShiftTypeRepositoryStub {
	return &ShiftTypeRepositoryStub{
		ShiftTypes: []*domain.ShiftType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ShiftTypeRepositoryStub) Query(dataModel *models.ShiftTypesQueryRequestParams) ([]*domain.ShiftType, error) {
	var ShiftTypes []*domain.ShiftType
	for _, v := range r.ShiftTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			ShiftTypes = append(ShiftTypes, v)
			break
		}
	}
	return ShiftTypes, nil
}

func (r *ShiftTypeRepositoryStub) Count(dataModel *models.ShiftTypesQueryRequestParams) (int64, error) {
	var ShiftTypes []*domain.ShiftType
	for _, v := range r.ShiftTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			ShiftTypes = append(ShiftTypes, v)
			break
		}
	}
	return int64(len(ShiftTypes)), nil
}

func (r *ShiftTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *ShiftTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *ShiftTypeRepositoryStub) GetShiftTypesByIds(ids []int64) ([]*domain.ShiftType, error) {
	var ShiftTypes []*domain.ShiftType
	for _, v := range r.ShiftTypes {
		for _, id := range ids {
			if v.ID == uint(id) {
				ShiftTypes = append(ShiftTypes, v)
			}
		}
	}
	return ShiftTypes, nil
}
