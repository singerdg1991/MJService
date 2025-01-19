package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/models"
)

type HolidayRepositoryStub struct {
	Holidays []*domain.Holiday
}

type holidayTestCondition struct {
	HasError bool
}

var UserTestCondition *holidayTestCondition = &holidayTestCondition{}

func NewHolidayRepositoryStub() *HolidayRepositoryStub {
	return &HolidayRepositoryStub{
		Holidays: []*domain.Holiday{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *HolidayRepositoryStub) Query(dataModel *models.HolidaysQueryRequestParams) ([]*domain.Holiday, error) {
	var holidays []*domain.Holiday
	for _, v := range r.Holidays {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			holidays = append(holidays, v)
			break
		}
	}
	return holidays, nil
}

func (r *HolidayRepositoryStub) Count(dataModel *models.HolidaysQueryRequestParams) (int64, error) {
	var holidays []*domain.Holiday
	for _, v := range r.Holidays {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			holidays = append(holidays, v)
			break
		}
	}
	return int64(len(holidays)), nil
}

func (r *HolidayRepositoryStub) Migrate() {
	// do stuff
}

func (r *HolidayRepositoryStub) Seed() {
	// do stuff
}

func (r *HolidayRepositoryStub) Create(payload *models.HolidaysCreateRequestBody) (*domain.Holiday, error) {
	panic("implement me")
}

func (r *HolidayRepositoryStub) Delete(payload *models.HolidaysDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *HolidayRepositoryStub) Update(payload *models.HolidaysCreateRequestBody, id int64) (*domain.Holiday, error) {
	panic("implement me")
}

func (r *HolidayRepositoryStub) GetHolidaysByIds(ids []int64) ([]*domain.Holiday, error) {
	panic("implement me")
}

func (r *HolidayRepositoryStub) UpdateStatus(payload *models.HolidaysUpdateStatusRequestBody, id int64) (*domain.Holiday, error) {
	panic("implement me")
}
