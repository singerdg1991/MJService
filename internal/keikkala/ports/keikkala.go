package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/keikkala/domain"
	"github.com/hoitek/Maja-Service/internal/keikkala/models"
)

type KeikkalaService interface {
	Query(dataModel *models.KeikkalasQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.KeikkalasCreateRequestBody) (*domain.Keikkala, error)
	Delete(payload *models.KeikkalasDeleteRequestBody) (*restypes.DeleteResponse, error)
	GetKeikkalaShiftsByIds(ids []int64) ([]*domain.Keikkala, error)
	FindByID(id int64) (*domain.Keikkala, error)
	QueryShiftStatistics(queries *models.KeikkalasQueryShiftStatisticsRequestParams) (*models.KeikkalasQueryShiftStatisticsResponseData, error)
}

type KeikkalaRepositoryPostgresDB interface {
	Query(dataModel *models.KeikkalasQueryRequestParams) ([]*domain.Keikkala, error)
	Count(dataModel *models.KeikkalasQueryRequestParams) (int64, error)
	Create(payload *models.KeikkalasCreateRequestBody) (*domain.Keikkala, error)
	Delete(payload *models.KeikkalasDeleteRequestBody) ([]int64, error)
	GetKeikkalaShiftsByIds(ids []int64) ([]*domain.Keikkala, error)
	QueryShiftStatistics(queries *models.KeikkalasQueryShiftStatisticsRequestParams) (int, int, int, error)
}
