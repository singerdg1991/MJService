package models

import "github.com/hoitek/Maja-Service/internal/staff/domain"

/*
 * @apiDefine: StaffsQueryChatMessagesResponseData
 */
type StaffsQueryChatMessagesResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []StaffsCreateChatMessageResponseData `json:"items" openapi:"$ref:StaffsCreateChatMessageResponseData;type:array"`
}

/*
 * @apiDefine: StaffsQueryChatMessagesResponse
 */
type StaffsQueryChatMessagesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryChatMessagesResponseData `json:"data" openapi:"$ref:StaffsQueryChatMessagesResponseData"`
}

/*
 * @apiDefine: StaffsQueryChatMessagesNotFoundResponse
 */
type StaffsQueryChatMessagesNotFoundResponse struct {
	Staffs []domain.StaffChatMessage `json:"staffs" openapi:"$ref:StaffChatMessage;type:array"`
}
