package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: TicketMessageTicket
 */
type TicketMessageTicket struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
}

/*
 * @apiDefine: TicketMessage
 */
type TicketMessage struct {
	ID            uint                    `json:"id" openapi:"example:1"`
	TicketID      *uint                   `json:"ticketId" openapi:"example:1"`
	SenderID      *uint                   `json:"senderId" openapi:"example:1"`
	RecipientID   *uint                   `json:"recipientId" openapi:"example:1"`
	Ticket        *TicketMessageTicket    `json:"ticket" openapi:"$ref:TicketMessageTicket;"`
	Sender        *TicketUser             `json:"sender" openapi:"$ref:TicketUser;"`
	Recipient     *TicketUser             `json:"recipient" openapi:"$ref:TicketUser;"`
	SenderType    string                  `json:"senderType" openapi:"example:customer"`
	RecipientType string                  `json:"recipientType" openapi:"example:customer"`
	Message       string                  `json:"message" openapi:"example:message;required"`
	Attachments   []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt     time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time               `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time              `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *TicketMessage) TableName() string {
	return "ticketMessage"
}

func (u *TicketMessage) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
