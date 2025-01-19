package models

import (
	"github.com/hoitek/Maja-Service/internal/ticket/domain"
)

/*
 * @apiDefine: TicketsQueryMessagesResponseData
 */
type TicketsQueryMessagesResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []TicketsCreateMessageResponseData `json:"items" openapi:"$ref:TicketsCreateMessageResponseData;type:array"`
}

/*
 * @apiDefine: TicketsQueryMessagesResponse
 */
type TicketsQueryMessagesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200;"`
	Data       TicketsQueryMessagesResponseData `json:"data" openapi:"$ref:TicketsQueryMessagesResponseData;type:object;"`
}

/*
 * @apiDefine: TicketsQueryMessagesNotFoundResponse
 */
type TicketsQueryMessagesNotFoundResponse struct {
	Tickets []domain.Ticket `json:"tickets" openapi:"$ref:Ticket;type:array"`
}
