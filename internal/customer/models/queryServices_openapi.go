package models

/*
 * @apiDefine: CustomersQueryServicesResponseData
 */
type CustomersQueryServicesResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreateServicesResponseData `json:"items" openapi:"$ref:CustomersCreateServicesResponseData"`
}

/*
 * @apiDefine: CustomersQueryServicesResponse
 */
type CustomersQueryServicesResponse struct {
	StatusCode int                                `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryServicesResponseData `json:"data" openapi:"$ref:CustomersQueryServicesResponseData"`
}

/*
 * @apiDefine: CustomersQueryServicesNotFoundResponse
 */
type CustomersQueryServicesNotFoundResponse struct {
	StatusCode int    `json:"statusCode" openapi:"example:404"`
	Message    string `json:"message" openapi:"example:Not Found"`
}
