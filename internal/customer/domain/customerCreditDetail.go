package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: CustomerCreditDetailBillingAddressStreet
 */
type CustomerCreditDetailBillingAddressStreet struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:streetName"`
}

/*
 * @apiDefine: CustomerCreditDetailBillingAddressCity
 */
type CustomerCreditDetailBillingAddressCity struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:cityName"`
}

/*
 * @apiDefine: CustomerCreditDetailBillingAddress
 */
type CustomerCreditDetailBillingAddress struct {
	ID             uint                                    `json:"id" openapi:"example:1"`
	City           *CustomerCreditDetailBillingAddressCity `json:"city" openapi:"$ref:CustomerCreditDetailBillingAddressCity"`
	Street         string                                  `json:"street" openapi:"example:streetName"`
	Name           string                                  `json:"name" openapi:"example:Home;required"`
	PostalCode     *string                                 `json:"postalCode" openapi:"example:1234567890;required"`
	BuildingNumber string                                  `json:"buildingNumber" openapi:"example:123;required"`
}

/*
 * @apiDefine: CustomerCreditDetail
 */
type CustomerCreditDetail struct {
	ID                uint                                `json:"id" openapi:"example:1"`
	CustomerID        uint                                `json:"customerId" openapi:"example:1"`
	BillingAddressID  uint                                `json:"billingAddressId" openapi:"example:1"`
	BillingAddress    *CustomerCreditDetailBillingAddress `json:"billingAddress" openapi:"$ref:CustomerCreditDetailBillingAddress"`
	BankAccountNumber string                              `json:"bankAccountNumber" openapi:"example:1234567890"`
	CreatedAt         time.Time                           `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt         time.Time                           `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt         *time.Time                          `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerCreditDetail() *CustomerCreditDetail {
	return &CustomerCreditDetail{}
}

func (ns *CustomerCreditDetail) TableName() string {
	return "customerCreditDetails"
}

func (ns *CustomerCreditDetail) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
