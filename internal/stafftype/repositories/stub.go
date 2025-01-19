package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/stafftype/domain"
	"github.com/hoitek/Maja-Service/internal/stafftype/models"
)

type StaffTypeRepositoryStub struct {
	StaffTypes []*domain.StaffType
}

type stafftypeTestCondition struct {
	HasError bool
}

var StaffTypeTestCondition *stafftypeTestCondition = &stafftypeTestCondition{}

func NewStaffTypeRepositoryStub() *StaffTypeRepositoryStub {
	return &StaffTypeRepositoryStub{
		StaffTypes: []*domain.StaffType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *StaffTypeRepositoryStub) Query(dataModel *models.StaffTypesQueryRequestParams) ([]*domain.StaffType, error) {
	var staffTypes []*domain.StaffType
	for _, v := range r.StaffTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			staffTypes = append(staffTypes, v)
			break
		}
	}
	return staffTypes, nil
}

func (r *StaffTypeRepositoryStub) Count(dataModel *models.StaffTypesQueryRequestParams) (int64, error) {
	var staffTypes []*domain.StaffType
	for _, v := range r.StaffTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			staffTypes = append(staffTypes, v)
			break
		}
	}
	return int64(len(staffTypes)), nil
}

func (r *StaffTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *StaffTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *StaffTypeRepositoryStub) GetStaffTypesByIds(ids []int64) ([]*domain.StaffType, error) {
	var staffTypes []*domain.StaffType
	for _, v := range r.StaffTypes {
		for _, id := range ids {
			if v.ID == uint(id) {
				staffTypes = append(staffTypes, v)
				break
			}
		}
	}
	return staffTypes, nil
}
