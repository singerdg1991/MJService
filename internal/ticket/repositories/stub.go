package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/ticket/domain"
	"github.com/hoitek/Maja-Service/internal/ticket/models"
)

type TicketRepositoryStub struct {
	Tickets []*domain.Ticket
}

type ticketCategoryTestCondition struct {
	HasError bool
}

var UserTestCondition *ticketCategoryTestCondition = &ticketCategoryTestCondition{}

func NewTicketRepositoryStub() *TicketRepositoryStub {
	return &TicketRepositoryStub{
		Tickets: []*domain.Ticket{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *TicketRepositoryStub) Query(dataModel *models.TicketsQueryRequestParams) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	for _, v := range r.Tickets {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			tickets = append(tickets, v)
			break
		}
	}
	return tickets, nil
}

func (r *TicketRepositoryStub) QueryMessages(payload *models.TicketsQueryMessagesRequestParams) ([]*domain.TicketMessage, error) {
	panic("implement me")
}

func (r *TicketRepositoryStub) Count(dataModel *models.TicketsQueryRequestParams) (int64, error) {
	var tickets []*domain.Ticket
	for _, v := range r.Tickets {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			tickets = append(tickets, v)
			break
		}
	}
	return int64(len(tickets)), nil
}

func (r *TicketRepositoryStub) CountMessages(queries *models.TicketsQueryMessagesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *TicketRepositoryStub) Create(payload *models.TicketsCreateRequestBody, createdBy uint, senderType string, recipientType string) (*domain.Ticket, error) {
	panic("implement me")
}

func (r *TicketRepositoryStub) CreateMessage(ticketID int64, payload *models.TicketsCreateMessageRequestBody, senderID int64, recipientId *int64, senderType string, recipientType string) (*domain.TicketMessage, error) {
	panic("implement me")
}

func (r *TicketRepositoryStub) Delete(payload *models.TicketsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}
