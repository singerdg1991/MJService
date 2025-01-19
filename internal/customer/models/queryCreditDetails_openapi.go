package models

import "github.com/hoitek/Maja-Service/internal/customer/domain"

/*
 * @apiDefine: CustomersQueryCreditDetailsResponseData
 */
type CustomersQueryCreditDetailsResponseData struct {
	Limit      int                                        `json:"limit" openapi:"example:10"`
	Offset     int                                        `json:"offset" openapi:"example:0"`
	Page       int                                        `json:"page" openapi:"example:1"`
	TotalRows  int                                        `json:"totalRows" openapi:"example:1"`
	TotalPages int                                        `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreateCreditDetailsResponseData `json:"items" openapi:"$ref:CustomersCreateCreditDetailsResponseData;type:array"`
}

/*
 * @apiDefine: CustomersQueryCreditDetailsResponse
 */
type CustomersQueryCreditDetailsResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryCreditDetailsResponseData `json:"data" openapi:"$ref:CustomersQueryCreditDetailsResponseData"`
}

/*
 * @apiDefine: CustomersQueryCreditDetailsNotFoundResponse
 */
type CustomersQueryCreditDetailsNotFoundResponse struct {
	Customers []domain.Customer `json:"customers" openapi:"$ref:Customer;type:array"`
}
