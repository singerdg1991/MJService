package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/service/domain"
	"github.com/hoitek/Maja-Service/internal/service/models"
)

type ServiceService interface {
	Query(dataModel *models.ServicesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ServicesCreateRequestBody) (*domain.Service, error)
	Delete(payload *models.ServicesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ServicesCreateRequestBody, id int64) (*domain.Service, error)
	GetServicesByIds(ids []int64) ([]*domain.Service, error)
	FindByID(id int64) (*domain.Service, error)
	FindByName(name string) (*domain.Service, error)
}

type ServiceRepositoryPostgresDB interface {
	Query(dataModel *models.ServicesQueryRequestParams) ([]*domain.Service, error)
	Count(dataModel *models.ServicesQueryRequestParams) (int64, error)
	Create(payload *models.ServicesCreateRequestBody) (*domain.Service, error)
	Delete(payload *models.ServicesDeleteRequestBody) ([]int64, error)
	Update(payload *models.ServicesCreateRequestBody, id int64) (*domain.Service, error)
	GetServicesByIds(ids []int64) ([]*domain.Service, error)
}
