package models

/*
 * @apiDefine: CustomersDeleteOtherAttachmentsResponseData
 */
type CustomersDeleteOtherAttachmentsResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteOtherAttachmentsResponse
 */
type CustomersDeleteOtherAttachmentsResponse struct {
	StatusCode int                                         `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteOtherAttachmentsResponseData `json:"data" openapi:"$ref:CustomersDeleteOtherAttachmentsResponseData"`
}
