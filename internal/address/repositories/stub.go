package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/address/domain"
	"github.com/hoitek/Maja-Service/internal/address/models"
)

type AddressRepositoryStub struct {
	Addresses []*domain.Address
}

type addressTestCondition struct {
	HasError bool
}

var UserTestCondition *addressTestCondition = &addressTestCondition{}

func NewAddressRepositoryStub() *AddressRepositoryStub {
	return &AddressRepositoryStub{
		Addresses: []*domain.Address{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *AddressRepositoryStub) Query(dataModel *models.AddressesQueryRequestParams) ([]*domain.Address, error) {
	var addresses []*domain.Address
	for _, v := range r.Addresses {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			addresses = append(addresses, v)
			break
		}
	}
	return addresses, nil
}

func (r *AddressRepositoryStub) Count(dataModel *models.AddressesQueryRequestParams) (int64, error) {
	var addresses []*domain.Address
	for _, v := range r.Addresses {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			addresses = append(addresses, v)
			break
		}
	}
	return int64(len(addresses)), nil
}

func (r *AddressRepositoryStub) Migrate() {
	// do stuff
}

func (r *AddressRepositoryStub) Seed() {
	// do stuff
}

func (r *AddressRepositoryStub) Create(payload *models.AddressesCreateRequestBody) (*domain.Address, error) {
	panic("implement me")
}

func (r *AddressRepositoryStub) Delete(payload *models.AddressesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *AddressRepositoryStub) Update(payload *models.AddressesCreateRequestBody, id int) (*domain.Address, error) {
	panic("implement me")
}
