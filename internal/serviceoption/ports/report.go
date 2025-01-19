package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/serviceoption/domain"
	"github.com/hoitek/Maja-Service/internal/serviceoption/models"
)

type ServiceOptionService interface {
	Query(dataModel *models.ServiceOptionsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ServiceOptionsCreateRequestBody) (*domain.ServiceOption, error)
	Delete(payload *models.ServiceOptionsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ServiceOptionsCreateRequestBody, id int64) (*domain.ServiceOption, error)
	GetServiceOptionsByIds(ids []int64) ([]*domain.ServiceOption, error)
	FindByID(id int64) (*domain.ServiceOption, error)
	FindByName(name string) (*domain.ServiceOption, error)
	FindByNameAndServiceTypeID(name string, serviceTypeId int) (*domain.ServiceOption, error)
}

type ServiceOptionRepositoryPostgresDB interface {
	Query(dataModel *models.ServiceOptionsQueryRequestParams) ([]*domain.ServiceOption, error)
	Count(dataModel *models.ServiceOptionsQueryRequestParams) (int64, error)
	Create(payload *models.ServiceOptionsCreateRequestBody) (*domain.ServiceOption, error)
	Delete(payload *models.ServiceOptionsDeleteRequestBody) ([]int64, error)
	Update(payload *models.ServiceOptionsCreateRequestBody, id int64) (*domain.ServiceOption, error)
	GetServiceOptionsByIds(ids []int64) ([]*domain.ServiceOption, error)
}
