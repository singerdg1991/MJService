package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/geartype/domain"
	"github.com/hoitek/Maja-Service/internal/geartype/models"
)

type GearTypeService interface {
	Query(dataModel *models.GearTypesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.GearTypesCreateRequestBody) (*domain.GearType, error)
	Delete(payload *models.GearTypesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.GearTypesCreateRequestBody, name string) (*domain.GearType, error)
}

type GearTypeRepositoryPostgresDB interface {
	Query(dataModel *models.GearTypesQueryRequestParams) ([]*domain.GearType, error)
	Count(dataModel *models.GearTypesQueryRequestParams) (int64, error)
	Create(payload *models.GearTypesCreateRequestBody) (*domain.GearType, error)
	Delete(payload *models.GearTypesDeleteRequestBody) ([]int64, error)
	Update(payload *models.GearTypesCreateRequestBody, name string) (*domain.GearType, error)
}
