package models

import "github.com/hoitek/Maja-Service/internal/cycle/domain"

/*
 * @apiDefine: CyclesQueryChatMessagesResponseData
 */
type CyclesQueryChatMessagesResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []CyclesCreateChatMessageResponseData `json:"items" openapi:"$ref:CyclesCreateChatMessageResponseData;type:array"`
}

/*
 * @apiDefine: CyclesQueryChatMessagesResponse
 */
type CyclesQueryChatMessagesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesQueryChatMessagesResponseData `json:"data" openapi:"$ref:CyclesQueryChatMessagesResponseData"`
}

/*
 * @apiDefine: CyclesQueryChatMessagesNotFoundResponse
 */
type CyclesQueryChatMessagesNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
