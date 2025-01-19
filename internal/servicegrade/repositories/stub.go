package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/servicegrade/domain"
	"github.com/hoitek/Maja-Service/internal/servicegrade/models"
)

type ServiceGradeRepositoryStub struct {
	ServiceGrades []*domain.ServiceGrade
}

type servicegradeTestCondition struct {
	HasError bool
}

var UserTestCondition *servicegradeTestCondition = &servicegradeTestCondition{}

func NewServiceGradeRepositoryStub() *ServiceGradeRepositoryStub {
	return &ServiceGradeRepositoryStub{
		ServiceGrades: []*domain.ServiceGrade{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ServiceGradeRepositoryStub) Query(dataModel *models.ServiceGradesQueryRequestParams) ([]*domain.ServiceGrade, error) {
	var servicegrades []*domain.ServiceGrade
	for _, v := range r.ServiceGrades {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			servicegrades = append(servicegrades, v)
			break
		}
	}
	return servicegrades, nil
}

func (r *ServiceGradeRepositoryStub) Count(dataModel *models.ServiceGradesQueryRequestParams) (int64, error) {
	var servicegrades []*domain.ServiceGrade
	for _, v := range r.ServiceGrades {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			servicegrades = append(servicegrades, v)
			break
		}
	}
	return int64(len(servicegrades)), nil
}

func (r *ServiceGradeRepositoryStub) Migrate() {
	// do stuff
}

func (r *ServiceGradeRepositoryStub) Seed() {
	// do stuff
}

func (r *ServiceGradeRepositoryStub) Create(payload *models.ServiceGradesCreateRequestBody) (*domain.ServiceGrade, error) {
	panic("implement me")
}

func (r *ServiceGradeRepositoryStub) Delete(payload *models.ServiceGradesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ServiceGradeRepositoryStub) Update(payload *models.ServiceGradesCreateRequestBody, id int64) (*domain.ServiceGrade, error) {
	panic("implement me")
}

func (r *ServiceGradeRepositoryStub) GetServiceGradesByIds(ids []int64) ([]*domain.ServiceGrade, error) {
	panic("implement me")
}
