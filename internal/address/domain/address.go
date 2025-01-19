package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: AddressStaff
 */
type AddressStaff struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	Email     string `json:"email" openapi:"example:email"`
}

/*
 * @apiDefine: AddressCustomer
 */
type AddressCustomer struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	Email     string `json:"email" openapi:"example:email"`
}

/*
 * @apiDefine: AddressCity
 */
type AddressCity struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:cityName"`
}

/*
 * @apiDefine: AddressStreet
 */
type AddressStreet struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:streetName"`
}

/*
 * @apiDefine: Address
 */
type Address struct {
	ID                uint             `json:"id" openapi:"example:1"`
	StaffID           *uint            `json:"-" openapi:"example:1"`
	CustomerID        *uint            `json:"-" openapi:"example:1"`
	CityID            uint             `json:"-" openapi:"example:1"`
	City              *AddressCity     `json:"city" openapi:"-"`
	Staff             *AddressStaff    `json:"staff" openapi:"-"`
	Customer          *AddressCustomer `json:"customer" openapi:"-"`
	Street            string           `json:"street" openapi:"example:streetName;required"`
	Name              string           `json:"name" openapi:"example:Home;required"`
	PostalCode        *string          `json:"postalCode" openapi:"example:1234567890;required"`
	BuildingNumber    string           `json:"buildingNumber" openapi:"example:123;required"`
	IsDeliveryAddress bool             `json:"isDeliveryAddress" openapi:"example:true;required"`
	IsMainAddress     bool             `json:"isMainAddress" openapi:"example:true;required"`
	CreatedAt         time.Time        `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt         time.Time        `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt         *time.Time       `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Address) TableName() string {
	return "addresses"
}

func (u *Address) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
