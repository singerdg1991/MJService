package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/address/domain"
	"github.com/hoitek/Maja-Service/internal/address/models"
)

type AddressService interface {
	Query(dataModel *models.AddressesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.AddressesCreateRequestBody) (*domain.Address, error)
	Delete(payload *models.AddressesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.AddressesCreateRequestBody, id int) (*domain.Address, error)
	FindByID(id int64) (*domain.Address, error)
}

type AddressRepositoryPostgresDB interface {
	Query(dataModel *models.AddressesQueryRequestParams) ([]*domain.Address, error)
	Count(dataModel *models.AddressesQueryRequestParams) (int64, error)
	Create(payload *models.AddressesCreateRequestBody) (*domain.Address, error)
	Delete(payload *models.AddressesDeleteRequestBody) ([]int64, error)
	Update(payload *models.AddressesCreateRequestBody, id int) (*domain.Address, error)
}
