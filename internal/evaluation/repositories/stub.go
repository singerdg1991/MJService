package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/evaluation/domain"
	"github.com/hoitek/Maja-Service/internal/evaluation/models"
)

type EvaluationRepositoryStub struct {
	Evaluations []*domain.Evaluation
}

type evaluationTestCondition struct {
	HasError bool
}

var UserTestCondition *evaluationTestCondition = &evaluationTestCondition{}

func NewEvaluationRepositoryStub() *EvaluationRepositoryStub {
	return &EvaluationRepositoryStub{
		Evaluations: []*domain.Evaluation{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *EvaluationRepositoryStub) Query(dataModel *models.EvaluationsQueryRequestParams) ([]*domain.Evaluation, error) {
	var evaluations []*domain.Evaluation
	for _, v := range r.Evaluations {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			evaluations = append(evaluations, v)
			break
		}
	}
	return evaluations, nil
}

func (r *EvaluationRepositoryStub) Count(dataModel *models.EvaluationsQueryRequestParams) (int64, error) {
	var evaluations []*domain.Evaluation
	for _, v := range r.Evaluations {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			evaluations = append(evaluations, v)
			break
		}
	}
	return int64(len(evaluations)), nil
}

func (r *EvaluationRepositoryStub) Migrate() {
	// do stuff
}

func (r *EvaluationRepositoryStub) Seed() {
	// do stuff
}

func (r *EvaluationRepositoryStub) Create(payload *models.EvaluationsCreateRequestBody) (*domain.Evaluation, error) {
	panic("implement me")
}

func (r *EvaluationRepositoryStub) Delete(payload *models.EvaluationsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *EvaluationRepositoryStub) Update(payload *models.EvaluationsCreateRequestBody, id int64) (*domain.Evaluation, error) {
	panic("implement me")
}

func (r *EvaluationRepositoryStub) GetEvaluationsByIds(ids []int64) ([]*domain.Evaluation, error) {
	panic("implement me")
}
