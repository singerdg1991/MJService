package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/punishment/domain"
	"github.com/hoitek/Maja-Service/internal/punishment/models"
)

type PunishmentService interface {
	Query(dataModel *models.PunishmentsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.PunishmentsCreateRequestBody) (*domain.Punishment, error)
	Delete(payload *models.PunishmentsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.PunishmentsCreateRequestBody, id int64) (*domain.Punishment, error)
	GetPunishmentsByIds(ids []int64) ([]*domain.Punishment, error)
	FindByID(id int64) (*domain.Punishment, error)
	FindByName(name string) (*domain.Punishment, error)
}

type PunishmentRepositoryPostgresDB interface {
	Query(dataModel *models.PunishmentsQueryRequestParams) ([]*domain.Punishment, error)
	Count(dataModel *models.PunishmentsQueryRequestParams) (int64, error)
	Create(payload *models.PunishmentsCreateRequestBody) (*domain.Punishment, error)
	Delete(payload *models.PunishmentsDeleteRequestBody) ([]int64, error)
	Update(payload *models.PunishmentsCreateRequestBody, id int64) (*domain.Punishment, error)
	GetPunishmentsByIds(ids []int64) ([]*domain.Punishment, error)
}
