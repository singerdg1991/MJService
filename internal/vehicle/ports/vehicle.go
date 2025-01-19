package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
	"github.com/hoitek/Maja-Service/internal/vehicle/models"
)

type VehicleService interface {
	Query(dataModel *models.VehiclesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.VehiclesCreateRequestBody) (*domain.Vehicle, error)
	Delete(payload *models.VehiclesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.VehiclesCreateRequestBody, id int) (*domain.Vehicle, error)
}

type VehicleRepositoryPostgresDB interface {
	Query(dataModel *models.VehiclesQueryRequestParams) ([]*domain.Vehicle, error)
	Count(dataModel *models.VehiclesQueryRequestParams) (int64, error)
	Create(payload *models.VehiclesCreateRequestBody) (*domain.Vehicle, error)
	Delete(payload *models.VehiclesDeleteRequestBody) ([]int64, error)
	Update(payload *models.VehiclesCreateRequestBody, id int) (*domain.Vehicle, error)
}
