package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/trash/domain"
	"github.com/hoitek/Maja-Service/internal/trash/models"
)

type TrashService interface {
	Query(dataModel *models.TrashesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.TrashesCreateRequestBody) (*domain.Trash, error)
	Delete(payload *models.TrashesDeleteRequestBody) (*restypes.DeleteResponse, error)
	FindByID(id int64) (*domain.Trash, error)
}

type TrashRepositoryPostgresDB interface {
	Query(dataModel *models.TrashesQueryRequestParams) ([]*domain.Trash, error)
	Count(dataModel *models.TrashesQueryRequestParams) (int64, error)
	Create(payload *models.TrashesCreateRequestBody) (*domain.Trash, error)
	Delete(payload *models.TrashesDeleteRequestBody) ([]int64, error)
}
