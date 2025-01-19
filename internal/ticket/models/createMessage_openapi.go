package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/ticket/domain"
)

/*
 * @apiDefine: TicketsCreateMessageResponseData
 */
type TicketsCreateMessageResponseData struct {
	ID            uint                       `json:"id" openapi:"example:1"`
	TicketID      uint                       `json:"ticketId" openapi:"example:1"`
	SenderId      uint                       `json:"senderId" openapi:"example:1"`
	RecipientId   uint                       `json:"recipientId" openapi:"example:1"`
	Ticket        domain.TicketMessageTicket `json:"ticket" openapi:"$ref:TicketMessageTicket;"`
	Sender        domain.TicketUser          `json:"sender" openapi:"$ref:TicketUser;"`
	Recipient     domain.TicketUser          `json:"recipient" openapi:"$ref:TicketUser;"`
	SenderType    string                     `json:"senderType" openapi:"example:customer"`
	RecipientType string                     `json:"recipientType" openapi:"example:customer"`
	Message       string                     `json:"message" openapi:"example:message;required"`
	Attachments   []types.UploadMetadata     `json:"attachments" openapi:"$ref:UploadMetadata;type:array"`
}

/*
 * @apiDefine: TicketsCreateMessageResponse
 */
type TicketsCreateMessageResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       TicketsCreateMessageResponseData `json:"data" openapi:"$ref:TicketsCreateMessageResponseData"`
}
