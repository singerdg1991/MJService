package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: CustomerRelativeCity
 */
type CustomerRelativeCity struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:Tehran"`
}

/*
 * @apiDefine: CustomerRelativeCustomer
 */
type CustomerRelativeCustomer struct {
	ID        uint   `json:"id" openapi:"example:1"`
	UserID    *uint  `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: CustomerRelative
 */
type CustomerRelative struct {
	ID                    uint                      `json:"id" openapi:"example:1"`
	CustomerID            uint                      `json:"customerId" openapi:"example:1"`
	Customer              *CustomerRelativeCustomer `json:"customer" openapi:"$ref:CustomerRelativeCustomer"`
	CityID                *uint                     `json:"cityId" openapi:"example:1"`
	City                  *CustomerRelativeCity     `json:"city" openapi:"$ref:CustomerRelativeCity"`
	FirstName             string                    `json:"firstName" openapi:"example:John;required"`
	LastName              string                    `json:"lastName" openapi:"example:Doe;ignored"`
	PhoneNumber           string                    `json:"phoneNumber" openapi:"example:09123456789"`
	Relation              string                    `json:"relation" openapi:"example:father"`
	AddressName           string                    `json:"addressName" openapi:"example:home"`
	AddressStreet         string                    `json:"addressStreet" openapi:"example:street"`
	AddressBuildingNumber string                    `json:"addressBuildingNumber" openapi:"example:1"`
	AddressPostalCode     string                    `json:"addressPostalCode" openapi:"example:1234567890"`
	CreatedAt             time.Time                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             time.Time                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             *time.Time                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerRelative() *CustomerRelative {
	return &CustomerRelative{}
}

func (ns *CustomerRelative) TableName() string {
	return "customerRelatives"
}

func (ns *CustomerRelative) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
