package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/limitation/domain"
	"github.com/hoitek/Maja-Service/internal/limitation/models"
)

type LimitationService interface {
	Query(dataModel *models.LimitationsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.LimitationsCreateRequestBody) (*domain.Limitation, error)
	Delete(payload *models.LimitationsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.LimitationsCreateRequestBody, id int64) (*domain.Limitation, error)
	GetLimitationsByIds(ids []int64) ([]*domain.Limitation, error)
	FindByID(id int64) (*domain.Limitation, error)
	FindByName(name string) (*domain.Limitation, error)
}

type LimitationRepositoryPostgresDB interface {
	Query(dataModel *models.LimitationsQueryRequestParams) ([]*domain.Limitation, error)
	Count(dataModel *models.LimitationsQueryRequestParams) (int64, error)
	Create(payload *models.LimitationsCreateRequestBody) (*domain.Limitation, error)
	Delete(payload *models.LimitationsDeleteRequestBody) ([]int64, error)
	Update(payload *models.LimitationsCreateRequestBody, id int64) (*domain.Limitation, error)
	GetLimitationsByIds(ids []int64) ([]*domain.Limitation, error)
}
