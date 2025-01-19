package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/ticket/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/ticket/domain"
	"github.com/hoitek/Maja-Service/internal/ticket/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type TicketRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewTicketRepositoryPostgresDB(d *sql.DB) *TicketRepositoryPostgresDB {
	return &TicketRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.TicketsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" t.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *TicketRepositoryPostgresDB) Query(queries *models.TicketsQueryRequestParams) ([]*domain.Ticket, error) {
	q := `SELECT * FROM tickets t `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.title %s", queries.Sorts.Title.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var tickets []*domain.Ticket
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ticket              domain.Ticket
			userID              sql.NullInt64
			departmentID        sql.NullInt64
			createdBy           sql.NullInt64
			updatedBy           sql.NullInt64
			deletedAt           sql.NullTime
			deletedBy           sql.NullInt64
			attachments         json.RawMessage
			attachmentsMetadata []*types.UploadMetadata
		)
		err := rows.Scan(
			&ticket.ID,
			&ticket.Code,
			&userID,
			&departmentID,
			&ticket.SenderType,
			&ticket.RecipientType,
			&ticket.Title,
			&ticket.Description,
			&ticket.Status,
			&ticket.Priority,
			&attachments,
			&ticket.CreatedAt,
			&createdBy,
			&ticket.UpdatedAt,
			&updatedBy,
			&deletedAt,
			&deletedBy,
		)
		if err != nil {
			return nil, err
		}
		if userID.Valid {
			ticket.UserID = uint(userID.Int64)
			ticket.User = &domain.TicketUser{}

			// Find user
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticket.UserID).Scan(
				&ticket.User.ID,
				&ticket.User.FirstName,
				&ticket.User.LastName,
				&ticket.User.Email,
				&ticket.User.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}
		if departmentID.Valid {
			ticket.DepartmentID = uint(departmentID.Int64)
		}
		if createdBy.Valid {
			ticket.CreatedByID = uint(createdBy.Int64)
			ticket.CreatedBy = domain.TicketUser{}
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticket.CreatedByID).Scan(
				&ticket.CreatedBy.ID,
				&ticket.CreatedBy.FirstName,
				&ticket.CreatedBy.LastName,
				&ticket.CreatedBy.Email,
				&ticket.CreatedBy.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}
		if updatedBy.Valid {
			uby := uint(updatedBy.Int64)
			ticket.UpdatedByID = &uby
			ticket.UpdatedBy = &domain.TicketUser{}
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticket.UpdatedByID).Scan(
				&ticket.UpdatedBy.ID,
				&ticket.UpdatedBy.FirstName,
				&ticket.UpdatedBy.LastName,
				&ticket.UpdatedBy.Email,
				&ticket.UpdatedBy.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}
		if deletedAt.Valid {
			ticket.DeletedAt = &deletedAt.Time
		}
		if deletedBy.Valid {
			dby := uint(deletedBy.Int64)
			ticket.DeletedByID = &dby
			ticket.DeletedBy = &domain.TicketUser{}
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticket.DeletedByID).Scan(
				&ticket.DeletedBy.ID,
				&ticket.DeletedBy.FirstName,
				&ticket.DeletedBy.LastName,
				&ticket.DeletedBy.Email,
				&ticket.DeletedBy.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}

		err = json.Unmarshal(attachments, &attachmentsMetadata)
		if err != nil {
			log.Printf("failed to unmarshal attachments metadata: %v in ticket: %d", err, ticket.ID)
		} else {
			for _, attachment := range attachmentsMetadata {
				attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.TICKET_BUCKET_NAME[len("maja."):])
			}
		}
		ticket.Attachments = attachmentsMetadata
		tickets = append(tickets, &ticket)
	}
	return tickets, nil
}

func (r *TicketRepositoryPostgresDB) Count(queries *models.TicketsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(t.id) FROM tickets t `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func makeWhereFiltersMessages(queries *models.TicketsQueryMessagesRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" tm.id = %d", queries.ID))
		}
		if queries.TicketID != 0 {
			where = append(where, fmt.Sprintf(" tm.ticketId = %d", queries.TicketID))
		}
	}
	return where
}

func (r *TicketRepositoryPostgresDB) QueryMessages(payload *models.TicketsQueryMessagesRequestParams) ([]*domain.TicketMessage, error) {
	q := `SELECT * FROM ticketMessages tm `
	if payload != nil {
		where := makeWhereFiltersMessages(payload)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if payload.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" tm.id %s", payload.Sorts.ID.Op))
		}
		if len(sorts) > 0 {
			q += " ORDER BY " + strings.Join(sorts, ",")
		}
		limit := exp.TerIf(payload.Limit == 0, 10, payload.Limit)
		payload.Page = exp.TerIf(payload.Page == 0, 1, payload.Page)
		offset := (payload.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	q += ";"

	var ticketMessages []*domain.TicketMessage
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ticketMessage domain.TicketMessage
			ticketID      sql.NullInt64
			senderID      sql.NullInt64
			recipientID   sql.NullInt64
			attachments   json.RawMessage
			deletedAt     sql.NullTime
		)
		err := rows.Scan(
			&ticketMessage.ID,
			&ticketID,
			&senderID,
			&recipientID,
			&ticketMessage.SenderType,
			&ticketMessage.RecipientType,
			&ticketMessage.Message,
			&attachments,
			&ticketMessage.CreatedAt,
			&ticketMessage.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if ticketID.Valid {
			ti := uint(ticketID.Int64)
			ticketMessage.TicketID = &ti
			ticketMessage.Ticket = &domain.TicketMessageTicket{}
			err := r.PostgresDB.QueryRow(`SELECT id, title FROM tickets WHERE id = $1`, ticketMessage.TicketID).Scan(
				&ticketMessage.Ticket.ID,
				&ticketMessage.Ticket.Title,
			)
			if err != nil {
				return nil, err
			}
		}
		if senderID.Valid {
			si := uint(senderID.Int64)
			ticketMessage.SenderID = &si
			ticketMessage.Sender = &domain.TicketUser{}
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticketMessage.SenderID).Scan(
				&ticketMessage.Sender.ID,
				&ticketMessage.Sender.FirstName,
				&ticketMessage.Sender.LastName,
				&ticketMessage.Sender.Email,
				&ticketMessage.Sender.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}
		if recipientID.Valid {
			ri := uint(recipientID.Int64)
			ticketMessage.RecipientID = &ri
			ticketMessage.Recipient = &domain.TicketUser{}
			err := r.PostgresDB.QueryRow(`SELECT id, firstName, lastName, email, avatarUrl FROM users WHERE id = $1`, ticketMessage.RecipientID).Scan(
				&ticketMessage.Recipient.ID,
				&ticketMessage.Recipient.FirstName,
				&ticketMessage.Recipient.LastName,
				&ticketMessage.Recipient.Email,
				&ticketMessage.Recipient.AvatarUrl,
			)
			if err != nil {
				return nil, err
			}
		}
		if deletedAt.Valid {
			ticketMessage.DeletedAt = &deletedAt.Time
		}
		err = json.Unmarshal(attachments, &ticketMessage.Attachments)
		if err != nil {
			log.Printf("failed to unmarshal attachments metadata: %v in ticket message: %d", err, ticketMessage.ID)
		} else {
			for _, attachment := range ticketMessage.Attachments {
				attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.TICKET_BUCKET_NAME[len("maja."):])
			}
		}
		ticketMessages = append(ticketMessages, &ticketMessage)
	}
	return ticketMessages, nil
}

func (r *TicketRepositoryPostgresDB) CountMessages(queries *models.TicketsQueryMessagesRequestParams) (int64, error) {
	q := `SELECT COUNT(tm.id) FROM ticketMessages tm `
	if queries != nil {
		where := makeWhereFiltersMessages(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TicketRepositoryPostgresDB) Create(payload *models.TicketsCreateRequestBody, createdBy uint, senderType string, recipientType string) (*domain.Ticket, error) {
	// Current time
	currentTime := time.Now()

	// Generate 8 length unique code
	code := fmt.Sprintf("#%d%d", time.Now().Unix(), createdBy)

	// Create fields dynamically
	var fields = map[string]interface{}{
		"code":          code,
		"title":         payload.Title,
		"description":   payload.Description,
		"status":        constants.TICKET_STATUS_OPEN,
		"priority":      constants.TICKET_PRIORITY_LOW,
		"createdBy":     createdBy,
		"senderType":    senderType,
		"recipientType": recipientType,
		"created_at":    currentTime,
		"updated_at":    currentTime,
	}
	if payload.UserID != nil {
		fields["userId"] = *payload.UserID
	}
	if payload.DepartmentID != nil {
		fields["departmentId"] = *payload.DepartmentID
	}

	// Create keys
	keys := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))
	for key, value := range fields {
		keys = append(keys, key)
		values = append(values, value)
	}

	// Prepare query
	q := `INSERT INTO tickets (` + strings.Join(keys, ",") + `) VALUES`
	indexes := make([]string, 0, len(keys))
	for i := range keys {
		indexes = append(indexes, fmt.Sprintf("$%d", i+1))
	}
	q += "(" + strings.Join(indexes, ",") + ") RETURNING id"
	log.Println(q)

	// Execute query
	var ticketID int64
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(q, values...).Scan(&ticketID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var (
		recipientId     *uint = nil
		ticketMessageID int64
	)
	if payload.UserID != nil {
		recipientId = payload.UserID
	}
	qMessage := `
		INSERT INTO ticketMessages (ticketId, senderId, recipientId, senderType, recipientType, message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err = tx.QueryRow(qMessage, ticketID, createdBy, recipientId, senderType, recipientType, payload.Description, currentTime, currentTime).Scan(&ticketMessageID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Get the ticketCategory
	tickets, err := r.Query(&models.TicketsQueryRequestParams{
		ID: int(ticketID),
	})
	if err != nil {
		return nil, err
	}
	if len(tickets) == 0 {
		return nil, errors.New("no ticket found")
	}
	ticket := *tickets[0]

	// Return the ticket
	return &ticket, nil
}

func (r *TicketRepositoryPostgresDB) CreateMessage(ticketID int64, payload *models.TicketsCreateMessageRequestBody, senderID int64, recipientId *int64, senderType string, recipientType string) (*domain.TicketMessage, error) {
	// Current time
	currentTime := time.Now()

	// Create fields dynamically
	var fields = map[string]interface{}{
		"ticketId":      ticketID,
		"senderId":      senderID,
		"recipientId":   recipientId,
		"senderType":    senderType,
		"recipientType": recipientType,
		"message":       payload.Message,
		"created_at":    currentTime,
		"updated_at":    currentTime,
		"deleted_at":    nil,
	}
	// Create keys
	keys := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))
	for key, value := range fields {
		keys = append(keys, key)
		values = append(values, value)
	}

	// Prepare query
	q := `INSERT INTO ticketMessages (` + strings.Join(keys, ",") + `) VALUES`
	indexes := make([]string, 0, len(keys))
	for i := range keys {
		indexes = append(indexes, fmt.Sprintf("$%d", i+1))
	}
	q += "(" + strings.Join(indexes, ",") + ") RETURNING id"
	log.Println(q)

	// Execute query
	var ticketMessageID int64
	err := r.PostgresDB.QueryRow(q, values...).Scan(&ticketMessageID)
	if err != nil {
		return nil, err
	}

	// Get the ticketCategory
	ticketMessages, err := r.QueryMessages(&models.TicketsQueryMessagesRequestParams{
		ID: int(ticketMessageID),
	})
	if err != nil {
		return nil, err
	}
	if len(ticketMessages) == 0 {
		return nil, errors.New("no ticket message found")
	}
	ticketMessage := *ticketMessages[0]

	// Return the ticket
	return &ticketMessage, nil
}

func (r *TicketRepositoryPostgresDB) Delete(payload *models.TicketsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM tickets
		WHERE id = ANY ($1)
		RETURNING id
	`, pq.Int64Array(payload.IDsInt64)).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return payload.IDsInt64, nil
}
