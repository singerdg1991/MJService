package models

import "github.com/hoitek/Maja-Service/internal/customer/domain"

/*
 * @apiDefine: CustomersQueryResponseData
 */
type CustomersQueryResponseData struct {
	Limit      int                                       `json:"limit" openapi:"example:10"`
	Offset     int                                       `json:"offset" openapi:"example:0"`
	Page       int                                       `json:"page" openapi:"example:1"`
	TotalRows  int                                       `json:"totalRows" openapi:"example:1"`
	TotalPages int                                       `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreatePersonalInfoResponseData `json:"items" openapi:"$ref:CustomersCreatePersonalInfoResponseData;type:array"`
}

/*
 * @apiDefine: CustomersQueryResponse
 */
type CustomersQueryResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryResponseData `json:"data" openapi:"$ref:CustomersQueryResponseData"`
}

/*
 * @apiDefine: CustomersQueryNotFoundResponse
 */
type CustomersQueryNotFoundResponse struct {
	Customers []domain.Customer `json:"customers" openapi:"$ref:Customer;type:array"`
}
