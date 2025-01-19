package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/ticket/constants"
	"github.com/hoitek/Maja-Service/internal/ticket/domain"
	"github.com/hoitek/Maja-Service/internal/ticket/models"
	"github.com/hoitek/Maja-Service/internal/ticket/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/exp"
)

type TicketService struct {
	PostgresRepository ports.TicketRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewTicketService(pDB ports.TicketRepositoryPostgresDB, m *storage.MinIO) TicketService {
	go minio.SetupMinIOStorage(constants.TICKET_BUCKET_NAME, m)
	return TicketService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *TicketService) Query(q *models.TicketsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying tickets", q)
	tickets, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.TicketsQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      tickets,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *TicketService) QueryMessages(q *models.TicketsQueryMessagesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying ticket messages", q)
	messages, err := s.PostgresRepository.QueryMessages(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountMessages(&models.TicketsQueryMessagesRequestParams{
		ID:       q.ID,
		TicketID: q.TicketID,
		Page:     q.Page,
		Limit:    0,
		Filters:  q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      messages,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *TicketService) Create(payload *models.TicketsCreateRequestBody, createdBy uint, senderType string, recipientType string) (*domain.Ticket, error) {
	return s.PostgresRepository.Create(payload, createdBy, senderType, recipientType)
}

func (s *TicketService) CreateMessage(ticketID int64, payload *models.TicketsCreateMessageRequestBody, senderID int64, recipientId *int64, senderType string, recipientType string) (*domain.TicketMessage, error) {
	return s.PostgresRepository.CreateMessage(ticketID, payload, senderID, recipientId, senderType, recipientType)
}

func (s *TicketService) Delete(payload *models.TicketsDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// TODO this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *TicketService) FindByID(id int64) (*domain.Ticket, error) {
	r, err := s.Query(&models.TicketsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("ticket not found")
	}
	tickets, ok := r.Items.([]*domain.Ticket)
	if !ok {
		return nil, errors.New("ticket not found")
	}
	if len(tickets) == 0 {
		return nil, errors.New("ticket not found")
	}
	return tickets[0], nil
}
