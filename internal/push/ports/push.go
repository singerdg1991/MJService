package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/push/domain"
	"github.com/hoitek/Maja-Service/internal/push/models"
)

type PushService interface {
	Query(dataModel *models.PushesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.PushesCreateRequestBody) (*domain.Push, error)
	GetPushesByIds(ids []int64) ([]*domain.Push, error)
	FindByID(id int64) (*domain.Push, error)
	FindByUserID(userID int) (*domain.Push, error)
}

type PushRepositoryPostgresDB interface {
	Query(dataModel *models.PushesQueryRequestParams) ([]*domain.Push, error)
	Count(dataModel *models.PushesQueryRequestParams) (int64, error)
	Create(payload *models.PushesCreateRequestBody) (*domain.Push, error)
	GetPushesByIds(ids []int64) ([]*domain.Push, error)
}
