package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/models"
)

type WarningService interface {
	Query(dataModel *models.WarningsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.WarningsCreateRequestBody) (*domain.Warning, error)
	Delete(payload *models.WarningsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.WarningsCreateRequestBody, id int64) (*domain.Warning, error)
	GetWarningsByIds(ids []int64) ([]*domain.Warning, error)
	FindByID(id int64) (*domain.Warning, error)
	FindByWarningNumber(warningNumber int64) (*domain.Warning, error)
}

type WarningRepositoryPostgresDB interface {
	Query(dataModel *models.WarningsQueryRequestParams) ([]*domain.Warning, error)
	Count(dataModel *models.WarningsQueryRequestParams) (int64, error)
	Create(payload *models.WarningsCreateRequestBody) (*domain.Warning, error)
	Delete(payload *models.WarningsDeleteRequestBody) ([]int64, error)
	Update(payload *models.WarningsCreateRequestBody, id int64) (*domain.Warning, error)
	GetWarningsByIds(ids []int64) ([]*domain.Warning, error)
}
