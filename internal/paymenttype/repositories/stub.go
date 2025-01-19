package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/paymenttype/domain"
	"github.com/hoitek/Maja-Service/internal/paymenttype/models"
)

type PaymentTypeRepositoryStub struct {
	PaymentTypes []*domain.PaymentType
}

type paymentTypeTestCondition struct {
	HasError bool
}

var UserTestCondition *paymentTypeTestCondition = &paymentTypeTestCondition{}

func NewPaymentTypeRepositoryStub() *PaymentTypeRepositoryStub {
	return &PaymentTypeRepositoryStub{
		PaymentTypes: []*domain.PaymentType{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *PaymentTypeRepositoryStub) Query(dataModel *models.PaymentTypesQueryRequestParams) ([]*domain.PaymentType, error) {
	var paymentTypes []*domain.PaymentType
	for _, v := range r.PaymentTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			paymentTypes = append(paymentTypes, v)
			break
		}
	}
	return paymentTypes, nil
}

func (r *PaymentTypeRepositoryStub) Count(dataModel *models.PaymentTypesQueryRequestParams) (int64, error) {
	var paymentTypes []*domain.PaymentType
	for _, v := range r.PaymentTypes {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			paymentTypes = append(paymentTypes, v)
			break
		}
	}
	return int64(len(paymentTypes)), nil
}

func (r *PaymentTypeRepositoryStub) Migrate() {
	// do stuff
}

func (r *PaymentTypeRepositoryStub) Seed() {
	// do stuff
}
