package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/email/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/email/domain"
	"github.com/hoitek/Maja-Service/internal/email/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type EmailRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewEmailRepositoryPostgresDB(d *sql.DB) *EmailRepositoryPostgresDB {
	return &EmailRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.EmailsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" e.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" e.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Email.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Email.Op, fmt.Sprintf("%v", queries.Filters.Email.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" e.email %s %s", opValue.Operator, val))
		}
		if queries.Filters.SenderID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.SenderID.Op, fmt.Sprintf("%v", queries.Filters.SenderID.Value))
			val := exp.TerIf(opValue.Value == "", "", opValue.Value)
			where = append(where, fmt.Sprintf(" e.senderId %s %s", opValue.Operator, val))
		}
		if queries.Filters.Category.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Category.Op, fmt.Sprintf("%v", queries.Filters.Category.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" e.category %s %s", opValue.Operator, val))
		}
		if queries.Filters.StarredAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StarredAt.Op, fmt.Sprintf("%v", queries.Filters.StarredAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" e.starred_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *EmailRepositoryPostgresDB) Query(queries *models.EmailsQueryRequestParams) ([]*domain.Email, error) {
	q := `
		SELECT
			e.id,
			e.senderId,
			e.email,
			e.cc,
			e.bcc,
			e.title,
			e.subject,
			e.message,
			e.attachments,
			e.category,
			e.starred_at,
			e.created_at,
			e.updated_at,
			e.deleted_at,
			u.id AS uID,
			u.firstName AS uFirstName,
			u.lastName AS uLastName,
			u.email AS uEmail,
			u.avatarUrl AS uAvatarUrl
		FROM emails e
		LEFT JOIN users u ON u.id = e.senderId
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" e.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" e.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" e.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var emails []*domain.Email
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			email         domain.Email
			senderId      sql.NullInt64
			cc            json.RawMessage
			bcc           json.RawMessage
			attachments   json.RawMessage
			starredAt     pq.NullTime
			deletedAt     pq.NullTime
			userId        sql.NullInt64
			userFirstName sql.NullString
			userLastName  sql.NullString
			userEmail     sql.NullString
			userAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&email.ID,
			&senderId,
			&email.ToEmail,
			&cc,
			&bcc,
			&email.Title,
			&email.Subject,
			&email.Message,
			&attachments,
			&email.Category,
			&starredAt,
			&email.CreatedAt,
			&email.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userEmail,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if senderId.Valid {
			sid := uint(senderId.Int64)
			email.SenderID = &sid
		}
		if cc != nil {
			var ccs []string
			err := json.Unmarshal(cc, &ccs)
			if err != nil {
				log.Printf("error unmarshalling cc: %v", err.Error())
			}
			email.Cc = ccs
		}
		if bcc != nil {
			var bccs []string
			err := json.Unmarshal(bcc, &bccs)
			if err != nil {
				log.Printf("error unmarshalling bcc: %v", err.Error())
			}
			email.Bcc = bccs
		}
		if attachments != nil {
			var attachmentsMetadata []*types.UploadMetadata
			err = json.Unmarshal(attachments, &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in email: %d", err, email.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.EMAIL_BUCKET_NAME[len("maja."):])
				}
			}
		}
		if starredAt.Valid {
			email.StarredAt = &starredAt.Time
		}
		if deletedAt.Valid {
			email.DeletedAt = &deletedAt.Time
		}
		if userId.Valid {
			email.Sender = &domain.EmailSender{
				ID:        uint(userId.Int64),
				FirstName: userFirstName.String,
				LastName:  userLastName.String,
				Email:     userEmail.String,
				AvatarUrl: userAvatarUrl.String,
			}
		}
		emails = append(emails, &email)
	}
	return emails, nil
}

func (r *EmailRepositoryPostgresDB) Count(queries *models.EmailsQueryRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(e.id)
		FROM emails e
		LEFT JOIN users u ON u.id = e.senderId
	`
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

func (r *EmailRepositoryPostgresDB) Create(payload *models.EmailsCreateRequestBody) (*domain.Email, error) {
	// Current time
	currentTime := time.Now()

	// Insert the email
	var insertedId int
	err := r.PostgresDB.QueryRow(`
		INSERT INTO emails (senderId, email, title, subject, message, category, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		payload.AuthenticatedUser.ID,
		payload.Email,
		payload.Title,
		payload.Subject,
		payload.Message,
		payload.Category,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the email
	emails, err := r.Query(&models.EmailsQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(emails) == 0 {
		return nil, errors.New("no rows affected")
	}
	return emails[0], nil
}

func (r *EmailRepositoryPostgresDB) Delete(payload *models.EmailsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM emails
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

func (r *EmailRepositoryPostgresDB) UpdateCategory(payload *models.EmailsUpdateCategoryRequestBody, id int64) (*domain.Email, error) {
	// Current time
	currentTime := time.Now()

	// Find the email by id
	foundEmails, _ := r.Query(&models.EmailsQueryRequestParams{ID: int(id)})
	if len(foundEmails) == 0 {
		return nil, errors.New("email not found")
	}
	foundEmail := foundEmails[0]
	if foundEmail.DeletedAt != nil {
		return nil, errors.New("email is deleted")
	}

	// Update the email
	var updatedId int
	err := r.PostgresDB.QueryRow(`
		UPDATE emails
		SET category = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`, payload.Category, currentTime, id).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the email
	emails, err := r.Query(&models.EmailsQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(emails) == 0 {
		return nil, errors.New("no rows affected")
	}
	return emails[0], nil
}

func (r *EmailRepositoryPostgresDB) UpdateStar(payload *models.EmailsUpdateStarRequestBody, id int64) (*domain.Email, error) {
	// Current time
	currentTime := time.Now()

	// Find the email by id
	foundEmails, _ := r.Query(&models.EmailsQueryRequestParams{ID: int(id)})
	if len(foundEmails) == 0 {
		return nil, errors.New("email not found")
	}
	foundEmail := foundEmails[0]
	if foundEmail.DeletedAt != nil {
		return nil, errors.New("email is deleted")
	}

	// Update the email
	var (
		updatedId int
		starredAt *time.Time = nil
	)
	if payload.IsStarredAsBool {
		starredAt = &currentTime
	}
	err := r.PostgresDB.QueryRow(`
		UPDATE emails
		SET starred_at = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`, starredAt, currentTime, id).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the email
	emails, err := r.Query(&models.EmailsQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(emails) == 0 {
		return nil, errors.New("no rows affected")
	}
	return emails[0], nil
}

func (r *EmailRepositoryPostgresDB) GetEmailsByIds(ids []int64) ([]*domain.Email, error) {
	var emails []*domain.Email
	rows, err := r.PostgresDB.Query(`
		SELECT
			e.id,
			e.senderId,
			e.email,
			e.cc,
			e.bcc,
			e.title,
			e.subject,
			e.message,
			e.attachments,
			e.category,
			e.starred_at,
			e.created_at,
			e.updated_at,
			e.deleted_at,
			u.id AS uID,
			u.firstName AS uFirstName,
			u.lastName AS uLastName,
			u.email AS uEmail,
			u.avatarUrl AS uAvatarUrl
		FROM emails e
		LEFT JOIN users u ON u.id = e.senderId
		WHERE e.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			email         domain.Email
			senderId      sql.NullInt64
			cc            json.RawMessage
			bcc           json.RawMessage
			attachments   json.RawMessage
			starredAt     pq.NullTime
			deletedAt     pq.NullTime
			userId        sql.NullInt64
			userFirstName sql.NullString
			userLastName  sql.NullString
			userEmail     sql.NullString
			userAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&email.ID,
			&senderId,
			&email.ToEmail,
			&cc,
			&bcc,
			&email.Title,
			&email.Subject,
			&email.Message,
			&attachments,
			&email.Category,
			&starredAt,
			&email.CreatedAt,
			&email.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userEmail,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if senderId.Valid {
			sid := uint(senderId.Int64)
			email.SenderID = &sid
		}
		if cc != nil {
			var ccs []string
			err := json.Unmarshal(cc, &ccs)
			if err != nil {
				log.Printf("error unmarshalling cc: %v", err.Error())
			}
			email.Cc = ccs
		}
		if bcc != nil {
			var bccs []string
			err := json.Unmarshal(bcc, &bccs)
			if err != nil {
				log.Printf("error unmarshalling bcc: %v", err.Error())
			}
			email.Bcc = bccs
		}
		if attachments != nil {
			var attachmentsMetadata []*types.UploadMetadata
			err = json.Unmarshal(attachments, &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in email: %d", err, email.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.EMAIL_BUCKET_NAME[len("maja."):])
				}
			}
		}
		if starredAt.Valid {
			email.StarredAt = &starredAt.Time
		}
		if deletedAt.Valid {
			email.DeletedAt = &deletedAt.Time
		}
		if userId.Valid {
			email.Sender = &domain.EmailSender{
				ID:        uint(userId.Int64),
				FirstName: userFirstName.String,
				LastName:  userLastName.String,
				Email:     userEmail.String,
				AvatarUrl: userAvatarUrl.String,
			}
		}
		emails = append(emails, &email)
	}
	return emails, nil
}
