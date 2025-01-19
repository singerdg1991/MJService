package models

import "github.com/hoitek/Maja-Service/internal/address/domain"

/*
 * @apiDefine: AddressesQueryResponseDataItem
 */
type AddressesQueryResponseDataItem struct {
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
 * @apiDefine: AddressesQueryResponseData
 */
type AddressesQueryResponseData struct {
	Limit      int                              `json:"limit" openapi:"example:10"`
	Offset     int                              `json:"offset" openapi:"example:0"`
	Page       int                              `json:"page" openapi:"example:1"`
	TotalRows  int                              `json:"totalRows" openapi:"example:1"`
	TotalPages int                              `json:"totalPages" openapi:"example:1"`
	Items      []AddressesQueryResponseDataItem `json:"items" openapi:"$ref:AddressesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: AddressesQueryResponse
 */
type AddressesQueryResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       AddressesQueryResponseData `json:"data" openapi:"$ref:AddressesQueryResponseData"`
}

/*
 * @apiDefine: AddressesQueryNotFoundResponse
 */
type AddressesQueryNotFoundResponse struct {
	Addresses []domain.Address `json:"addresses" openapi:"$ref:Address;type:array"`
}
