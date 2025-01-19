package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/contracttype/domain"
	"github.com/hoitek/Maja-Service/internal/contracttype/models"
)

type ContractTypeService interface {
	Query(dataModel *models.ContractTypesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ContractTypesCreateRequestBody) (*domain.ContractType, error)
	Delete(payload *models.ContractTypesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ContractTypesCreateRequestBody, id int64) (*domain.ContractType, error)
	GetContractTypesByIds(ids []int64) ([]*domain.ContractType, error)
	FindByID(id int64) (*domain.ContractType, error)
	FindByName(name string) (*domain.ContractType, error)
}

type ContractTypeRepositoryPostgresDB interface {
	Query(dataModel *models.ContractTypesQueryRequestParams) ([]*domain.ContractType, error)
	Count(dataModel *models.ContractTypesQueryRequestParams) (int64, error)
	Create(payload *models.ContractTypesCreateRequestBody) (*domain.ContractType, error)
	Delete(payload *models.ContractTypesDeleteRequestBody) ([]int64, error)
	Update(payload *models.ContractTypesCreateRequestBody, id int64) (*domain.ContractType, error)
	GetContractTypesByIds(ids []int64) ([]*domain.ContractType, error)
}
