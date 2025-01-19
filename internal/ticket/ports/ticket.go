package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/ticket/domain"
	"github.com/hoitek/Maja-Service/internal/ticket/models"
)

type TicketService interface {
	Query(dataModel *models.TicketsQueryRequestParams) (*restypes.QueryResponse, error)
	QueryMessages(dataModel *models.TicketsQueryMessagesRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.TicketsCreateRequestBody, createdBy uint, senderType string, recipientType string) (*domain.Ticket, error)
	CreateMessage(ticketID int64, payload *models.TicketsCreateMessageRequestBody, senderID int64, recipientId *int64, senderType string, recipientType string) (*domain.TicketMessage, error)
	Delete(payload *models.TicketsDeleteRequestBody) (*restypes.DeleteResponse, error)
	FindByID(id int64) (*domain.Ticket, error)
}

type TicketRepositoryPostgresDB interface {
	Query(dataModel *models.TicketsQueryRequestParams) ([]*domain.Ticket, error)
	QueryMessages(dataModel *models.TicketsQueryMessagesRequestParams) ([]*domain.TicketMessage, error)
	Count(dataModel *models.TicketsQueryRequestParams) (int64, error)
	CountMessages(dataModel *models.TicketsQueryMessagesRequestParams) (int64, error)
	Create(payload *models.TicketsCreateRequestBody, createdBy uint, senderType string, recipientType string) (*domain.Ticket, error)
	CreateMessage(ticketID int64, payload *models.TicketsCreateMessageRequestBody, senderID int64, recipientId *int64, senderType string, recipientType string) (*domain.TicketMessage, error)
	Delete(payload *models.TicketsDeleteRequestBody) ([]int64, error)
}
