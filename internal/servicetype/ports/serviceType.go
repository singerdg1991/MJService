package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/servicetype/domain"
	"github.com/hoitek/Maja-Service/internal/servicetype/models"
)

type ServiceTypeService interface {
	Query(dataModel *models.ServiceTypesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ServiceTypesCreateRequestBody) (*domain.ServiceType, error)
	Delete(payload *models.ServiceTypesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ServiceTypesCreateRequestBody, id int64) (*domain.ServiceType, error)
	GetServiceTypesByIds(ids []int64) ([]*domain.ServiceType, error)
	FindByID(id int64) (*domain.ServiceType, error)
	FindByName(name string) (*domain.ServiceType, error)
	FindByNameAndServiceID(name string, serviceId int) (*domain.ServiceType, error)
}

type ServiceTypeRepositoryPostgresDB interface {
	Query(dataModel *models.ServiceTypesQueryRequestParams) ([]*domain.ServiceType, error)
	Count(dataModel *models.ServiceTypesQueryRequestParams) (int64, error)
	Create(payload *models.ServiceTypesCreateRequestBody) (*domain.ServiceType, error)
	Delete(payload *models.ServiceTypesDeleteRequestBody) ([]int64, error)
	Update(payload *models.ServiceTypesCreateRequestBody, id int64) (*domain.ServiceType, error)
	GetServiceTypesByIds(ids []int64) ([]*domain.ServiceType, error)
}
