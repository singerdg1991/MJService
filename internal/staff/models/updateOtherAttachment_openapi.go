package models

/*
 * @apiDefine: StaffsUpdateOtherAttachmentResponse
 */
type StaffsUpdateOtherAttachmentResponse struct {
	StatusCode int                                      `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateOtherAttachmentsResponseData `json:"data" openapi:"$ref:StaffsCreateOtherAttachmentsResponseData"`
}
