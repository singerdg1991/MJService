package models

/*
 * @apiDefine: StaffsQueryOtherAttachmentsResponseData
 */
type StaffsQueryOtherAttachmentsResponseData struct {
	Limit      int                                        `json:"limit" openapi:"example:10"`
	Offset     int                                        `json:"offset" openapi:"example:0"`
	Page       int                                        `json:"page" openapi:"example:1"`
	TotalRows  int                                        `json:"totalRows" openapi:"example:1"`
	TotalPages int                                        `json:"totalPages" openapi:"example:1"`
	Items      []StaffsCreateOtherAttachmentsResponseData `json:"items" openapi:"$ref:StaffsCreateOtherAttachmentsResponseData;type:array;required"`
}

/*
 * @apiDefine: StaffsQueryOtherAttachmentsResponse
 */
type StaffsQueryOtherAttachmentsResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryOtherAttachmentsResponseData `json:"data" openapi:"$ref:StaffsQueryOtherAttachmentsResponseData"`
}
