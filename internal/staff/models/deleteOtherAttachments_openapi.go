package models

/*
 * @apiDefine: StaffsDeleteOtherAttachmentsResponseData
 */
type StaffsDeleteOtherAttachmentsResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: StaffsDeleteOtherAttachmentsResponse
 */
type StaffsDeleteOtherAttachmentsResponse struct {
	StatusCode int                                      `json:"statusCode" openapi:"example:200"`
	Data       StaffsDeleteOtherAttachmentsResponseData `json:"data" openapi:"$ref:StaffsDeleteOtherAttachmentsResponseData"`
}
