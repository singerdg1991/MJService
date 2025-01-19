package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/vehicletype/domain"
	"github.com/hoitek/Maja-Service/internal/vehicletype/models"
)

type VehicleTypeService interface {
	Query(dataModel *models.VehicleTypesQueryRequestParams) (*restypes.QueryResponse, error)
	FindByName(name string) (*domain.VehicleType, error)
	Create(payload *models.VehicleTypesCreateRequestBody) (*domain.VehicleType, error)
	Delete(payload *models.VehicleTypesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.VehicleTypesCreateRequestBody, name string) (*domain.VehicleType, error)
}

type VehicleTypeRepositoryPostgresDB interface {
	Query(dataModel *models.VehicleTypesQueryRequestParams) ([]*domain.VehicleType, error)
	Count(dataModel *models.VehicleTypesQueryRequestParams) (int64, error)
	Create(payload *models.VehicleTypesCreateRequestBody) (*domain.VehicleType, error)
	Delete(payload *models.VehicleTypesDeleteRequestBody) ([]int64, error)
	Update(payload *models.VehicleTypesCreateRequestBody, name string) (*domain.VehicleType, error)
}
