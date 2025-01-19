package models

import (
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreateRelativesResponseData
 */
type CustomersCreateRelativesResponseData struct {
	ID                    uint                             `json:"id" openapi:"example:1"`
	CustomerID            uint                             `json:"customerId" openapi:"example:1"`
	Customer              *domain.CustomerRelativeCustomer `json:"customer" openapi:"$ref:CustomerRelativeCustomer"`
	CityID                *uint                            `json:"cityId" openapi:"example:1"`
	City                  *domain.CustomerRelativeCity     `json:"city" openapi:"$ref:CustomerRelativeCity"`
	FirstName             string                           `json:"firstName" openapi:"example:John;required"`
	LastName              string                           `json:"lastName" openapi:"example:Doe;ignored"`
	PhoneNumber           string                           `json:"phoneNumber" openapi:"example:09123456789"`
	Relation              string                           `json:"relation" openapi:"example:father"`
	AddressName           string                           `json:"addressName" openapi:"example:home"`
	AddressStreet         string                           `json:"addressStreet" openapi:"example:street"`
	AddressBuildingNumber string                           `json:"addressBuildingNumber" openapi:"example:1"`
	AddressPostalCode     string                           `json:"addressPostalCode" openapi:"example:1234567890"`
	CreatedAt             string                           `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             string                           `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             *string                          `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersCreateRelativesResponse
 */
type CustomersCreateRelativesResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateRelativesResponseData `json:"data" openapi:"$ref:CustomersCreateRelativesResponseData"`
}
