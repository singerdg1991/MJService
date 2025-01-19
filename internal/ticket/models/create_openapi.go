package models

/*
 * @apiDefine: TicketsResponseData
 */
type TicketsResponseData struct {
	ID               uint   `json:"id" openapi:"example:1"`
	TicketCategoryID uint   `json:"ticketCategoryId" openapi:"example:1"`
	Name             string `json:"name" openapi:"example:saeed"`
	Description      string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: TicketsCreateResponse
 */
type TicketsCreateResponse struct {
	StatusCode int                 `json:"statusCode" openapi:"example:200"`
	Data       TicketsResponseData `json:"data" openapi:"$ref:TicketsResponseData"`
}
