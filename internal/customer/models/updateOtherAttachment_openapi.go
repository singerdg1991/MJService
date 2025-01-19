package models

/*
 * @apiDefine: CustomersUpdateOtherAttachmentResponse
 */
type CustomersUpdateOtherAttachmentResponse struct {
	StatusCode int                                         `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateOtherAttachmentsResponseData `json:"data" openapi:"$ref:CustomersCreateOtherAttachmentsResponseData"`
}
