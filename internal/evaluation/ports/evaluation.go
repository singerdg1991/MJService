package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/evaluation/domain"
	"github.com/hoitek/Maja-Service/internal/evaluation/models"
)

type EvaluationService interface {
	Query(dataModel *models.EvaluationsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.EvaluationsCreateRequestBody) (*domain.Evaluation, error)
	Delete(payload *models.EvaluationsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.EvaluationsCreateRequestBody, id int64) (*domain.Evaluation, error)
	GetEvaluationsByIds(ids []int64) ([]*domain.Evaluation, error)
	FindByID(id int64) (*domain.Evaluation, error)
	FindByTitle(title string) (*domain.Evaluation, error)
}

type EvaluationRepositoryPostgresDB interface {
	Query(dataModel *models.EvaluationsQueryRequestParams) ([]*domain.Evaluation, error)
	Count(dataModel *models.EvaluationsQueryRequestParams) (int64, error)
	Create(payload *models.EvaluationsCreateRequestBody) (*domain.Evaluation, error)
	Delete(payload *models.EvaluationsDeleteRequestBody) ([]int64, error)
	Update(payload *models.EvaluationsCreateRequestBody, id int64) (*domain.Evaluation, error)
	GetEvaluationsByIds(ids []int64) ([]*domain.Evaluation, error)
}
