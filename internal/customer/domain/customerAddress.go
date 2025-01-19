package domain

import "time"

/*
 * @apiDefine: CustomerAddressAddress
 */
type CustomerAddressAddress struct {
	Name    string `json:"name" openapi:"example:Home"`
	City    string `json:"city" openapi:"example:tehran"`
	ZipCode string `json:"zipCode" openapi:"example:1234567890"`
	State   string `json:"state" openapi:"example:tehran"`
}

/*
 * @apiDefine: CustomerAddress
 */
type CustomerAddress struct {
	ID         uint                   `json:"id" openapi:"example:1"`
	CustomerID uint                   `json:"customerId" openapi:"example:1"`
	AddressID  uint                   `json:"addressId" openapi:"example:1"`
	Address    CustomerAddressAddress `json:"address" openapi:"$ref:CustomerAddressAddress"`
	CreatedAt  time.Time              `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt  time.Time              `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt  *time.Time             `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}
