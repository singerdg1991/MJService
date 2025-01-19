package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/equipment/domain"
	"github.com/hoitek/Maja-Service/internal/equipment/models"
)

type EquipmentService interface {
	Query(dataModel *models.EquipmentsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.EquipmentsCreateRequestBody) (*domain.Equipment, error)
	Delete(payload *models.EquipmentsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.EquipmentsCreateRequestBody, id int64) (*domain.Equipment, error)
	GetEquipmentsByIds(ids []int64) ([]*domain.Equipment, error)
	FindByID(id int64) (*domain.Equipment, error)
	FindByName(name string) (*domain.Equipment, error)
}

type EquipmentRepositoryPostgresDB interface {
	Query(dataModel *models.EquipmentsQueryRequestParams) ([]*domain.Equipment, error)
	Count(dataModel *models.EquipmentsQueryRequestParams) (int64, error)
	Create(payload *models.EquipmentsCreateRequestBody) (*domain.Equipment, error)
	Delete(payload *models.EquipmentsDeleteRequestBody) ([]int64, error)
	Update(payload *models.EquipmentsCreateRequestBody, id int64) (*domain.Equipment, error)
	GetEquipmentsByIds(ids []int64) ([]*domain.Equipment, error)
}
