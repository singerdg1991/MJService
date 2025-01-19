package models

/*
 * @apiDefine: CustomersQueryRelativesResponseData
 */
type CustomersQueryRelativesResponseData struct {
	Limit      int                                    `json:"limit" openapi:"example:10"`
	Offset     int                                    `json:"offset" openapi:"example:0"`
	Page       int                                    `json:"page" openapi:"example:1"`
	TotalRows  int                                    `json:"totalRows" openapi:"example:1"`
	TotalPages int                                    `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreateRelativesResponseData `json:"items" openapi:"$ref:CustomersCreateRelativesResponseData;type:array"`
}

/*
 * @apiDefine: CustomersQueryRelativesResponse
 */
type CustomersQueryRelativesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryRelativesResponseData `json:"data" openapi:"$ref:CustomersQueryRelativesResponseData"`
}
