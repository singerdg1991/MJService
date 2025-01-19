package models

import "github.com/hoitek/Maja-Service/internal/address/domain"

/*
 * @apiDefine: AddressesCreateResponseData
 */
type AddressesCreateResponseData struct {
	ID                int                    `json:"id" openapi:"example:1"`
	City              domain.AddressCity     `json:"city" openapi:"$ref:AddressCity;type:object;"`
	Staff             domain.AddressStaff    `json:"staff" openapi:"$ref:AddressStaff;type:object;"`
	Customer          domain.AddressCustomer `json:"customer" openapi:"$ref:AddressCustomer;type:object;"`
	Street            string                 `json:"street" openapi:"example:streetName"`
	Name              string                 `json:"name" openapi:"example:test"`
	PostalCode        string                 `json:"postalCode" openapi:"example:1234567890"`
	BuildingNumber    string                 `json:"buildingNumber" openapi:"example:1D34"`
	IsDeliveryAddress bool                   `json:"isDeliveryAddress" openapi:"example:true"`
	IsMainAddress     bool                   `json:"isMainAddress" openapi:"example:true"`
	CreatedAt         string                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt         string                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt         *string                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z;nullable"`
}

/*
 * @apiDefine: AddressesCreateResponse
 */
type AddressesCreateResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200"`
	Data       AddressesCreateResponseData `json:"data" openapi:"$ref:AddressesCreateResponseData"`
}
