package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/diagnose/domain"
	"github.com/hoitek/Maja-Service/internal/diagnose/models"
)

type DiagnoseRepositoryStub struct {
	Diagnoses []*domain.Diagnose
}

type diagnoseTestCondition struct {
	HasError bool
}

var UserTestCondition *diagnoseTestCondition = &diagnoseTestCondition{}

func NewDiagnoseRepositoryStub() *DiagnoseRepositoryStub {
	return &DiagnoseRepositoryStub{
		Diagnoses: []*domain.Diagnose{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *DiagnoseRepositoryStub) Query(dataModel *models.DiagnosesQueryRequestParams) ([]*domain.Diagnose, error) {
	var diagnoses []*domain.Diagnose
	for _, v := range r.Diagnoses {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			diagnoses = append(diagnoses, v)
			break
		}
	}
	return diagnoses, nil
}

func (r *DiagnoseRepositoryStub) Count(dataModel *models.DiagnosesQueryRequestParams) (int64, error) {
	var diagnoses []*domain.Diagnose
	for _, v := range r.Diagnoses {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			diagnoses = append(diagnoses, v)
			break
		}
	}
	return int64(len(diagnoses)), nil
}

func (r *DiagnoseRepositoryStub) Migrate() {
	// do stuff
}

func (r *DiagnoseRepositoryStub) Seed() {
	// do stuff
}

func (r *DiagnoseRepositoryStub) Create(payload *models.DiagnosesCreateRequestBody) (*domain.Diagnose, error) {
	panic("implement me")
}

func (r *DiagnoseRepositoryStub) Delete(payload *models.DiagnosesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *DiagnoseRepositoryStub) Update(payload *models.DiagnosesCreateRequestBody, id int) (*domain.Diagnose, error) {
	panic("implement me")
}
