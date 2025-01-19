package models

/*
 * @apiDefine: CustomersQueryOtherAttachmentsResponseData
 */
type CustomersQueryOtherAttachmentsResponseData struct {
	Limit      int                                           `json:"limit" openapi:"example:10"`
	Offset     int                                           `json:"offset" openapi:"example:0"`
	Page       int                                           `json:"page" openapi:"example:1"`
	TotalRows  int                                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                                           `json:"totalPages" openapi:"example:1"`
	Items      []CustomersCreateOtherAttachmentsResponseData `json:"items" openapi:"$ref:CustomersCreateOtherAttachmentsResponseData;type:array;required"`
}

/*
 * @apiDefine: CustomersQueryOtherAttachmentsResponse
 */
type CustomersQueryOtherAttachmentsResponse struct {
	StatusCode int                                        `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryOtherAttachmentsResponseData `json:"data" openapi:"$ref:CustomersQueryOtherAttachmentsResponseData"`
}
