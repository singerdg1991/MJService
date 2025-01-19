package models

import "github.com/hoitek/Maja-Service/internal/ticket/domain"

/*
 * @apiDefine: TicketsQueryResponseDataItem
 */
type TicketsQueryResponseDataItem struct {
	ID               uint   `json:"id" openapi:"example:1"`
	TicketCategoryID uint   `json:"ticketCategoryId" openapi:"example:1"`
	Name             string `json:"name" openapi:"example:John;required"`
	Description      string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: TicketsQueryResponseData
 */
type TicketsQueryResponseData struct {
	Limit      int                            `json:"limit" openapi:"example:10"`
	Offset     int                            `json:"offset" openapi:"example:0"`
	Page       int                            `json:"page" openapi:"example:1"`
	TotalRows  int                            `json:"totalRows" openapi:"example:1"`
	TotalPages int                            `json:"totalPages" openapi:"example:1"`
	Items      []TicketsQueryResponseDataItem `json:"items" openapi:"$ref:TicketsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: TicketsQueryResponse
 */
type TicketsQueryResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200;"`
	Data       TicketsQueryResponseData `json:"data" openapi:"$ref:TicketsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: TicketsQueryNotFoundResponse
 */
type TicketsQueryNotFoundResponse struct {
	Tickets []domain.Ticket `json:"tickets" openapi:"$ref:Ticket;type:array"`
}
