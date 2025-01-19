package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/contracttype/domain"
	"github.com/hoitek/Maja-Service/internal/contracttype/models"
)

type ContractTypeRepositoryStub struct {
	ContractTypes []*domain.ContractType
}

type contractTypeTestCondition struct {
	HasError bool
}

var UserTestCondition *contractTypeTestCondition = &contractTypeTestCondition{}

func NewContractTypeRepositoryStub() *ContractTypeRepositoryStub {
	return &ContractTypeRepositoryStub{
		ContractTypes: []*domain.ContractType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *ContractTypeRepositoryStub) Query(dataModel *models.ContractTypesQueryRequestParams) ([]*domain.ContractType, error) {
	var contractTypes []*domain.ContractType
	for _, v := range r.ContractTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			contractTypes = append(contractTypes, v)
			break
		}
	}
	return contractTypes, nil
}

func (r *ContractTypeRepositoryStub) Count(dataModel *models.ContractTypesQueryRequestParams) (int64, error) {
	var contractTypes []*domain.ContractType
	for _, v := range r.ContractTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			contractTypes = append(contractTypes, v)
			break
		}
	}
	return int64(len(contractTypes)), nil
}

func (r *ContractTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *ContractTypeRepositoryStub) Seed() {
	// do stuff
}

func (r *ContractTypeRepositoryStub) Create(payload *models.ContractTypesCreateRequestBody) (*domain.ContractType, error) {
	panic("implement me")
}

func (r *ContractTypeRepositoryStub) Delete(payload *models.ContractTypesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ContractTypeRepositoryStub) Update(payload *models.ContractTypesCreateRequestBody, id int64) (*domain.ContractType, error) {
	panic("implement me")
}

func (r *ContractTypeRepositoryStub) GetContractTypesByIds(ids []int64) ([]*domain.ContractType, error) {
	panic("implement me")
}
