package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/prescription/constants"
	"github.com/hoitek/Maja-Service/internal/prescription/domain"
	"github.com/hoitek/Maja-Service/internal/prescription/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type PrescriptionRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewPrescriptionRepositoryPostgresDB(d *sql.DB) *PrescriptionRepositoryPostgresDB {
	return &PrescriptionRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.PrescriptionsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" p.id = %d", queries.ID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf(" p.customerId = %d", queries.CustomerID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" p.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *PrescriptionRepositoryPostgresDB) Query(queries *models.PrescriptionsQueryRequestParams) ([]*domain.Prescription, error) {
	q := `
		SELECT
			p.id,
			p.customerId,
			p.title,
			p.datetime,
			p.doctorFullName,
            p.start_date,
            p.end_date,
			p.status,
			p.attachments,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			c.id AS cID,
			u.id AS uID,
			u.firstName AS uFirstName,
			u.lastName AS uLastName,
			u.avatarUrl AS uAvatarUrl
		FROM prescriptions p
		LEFT JOIN customers c ON p.customerId = c.id
		LEFT JOIN users u ON c.userId = u.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" p.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" p.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var prescriptions []*domain.Prescription
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			prescription        domain.Prescription
			customerId          sql.NullInt64
			datetime            pq.NullTime
			startDate           pq.NullTime
			endDate             pq.NullTime
			attachments         sql.NullString
			attachmentsMetadata []*types.UploadMetadata
			deletedAt           pq.NullTime
			cID                 sql.NullInt64
			uID                 sql.NullInt64
			uFirstName          sql.NullString
			uLastName           sql.NullString
			uAvatarUrl          sql.NullString
		)
		err := rows.Scan(
			&prescription.ID,
			&customerId,
			&prescription.Title,
			&datetime,
			&prescription.DoctorFullName,
			&startDate,
			&endDate,
			&prescription.Status,
			&attachments,
			&prescription.CreatedAt,
			&prescription.UpdatedAt,
			&deletedAt,
			&cID,
			&uID,
			&uFirstName,
			&uLastName,
			&uAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if customerId.Valid {
			cid := uint(customerId.Int64)
			prescription.CustomerID = cid
		}
		if datetime.Valid {
			prescription.DateTime = &datetime.Time
		}
		if startDate.Valid {
			prescription.StartDate = &startDate.Time
		}
		if endDate.Valid {
			prescription.EndDate = &endDate.Time
		}
		if deletedAt.Valid {
			prescription.DeletedAt = &deletedAt.Time
		}
		if cID.Valid {
			prescription.Customer = &domain.PrescriptionCustomer{
				ID: cID.Int64,
			}
			if uID.Valid {
				prescription.Customer.UserID = uID.Int64
			}
			if uFirstName.Valid {
				prescription.Customer.FirstName = uFirstName.String
			}
			if uLastName.Valid {
				prescription.Customer.LastName = uLastName.String
			}
			if uAvatarUrl.Valid {
				prescription.Customer.AvatarUrl = uAvatarUrl.String
			}
		}
		if attachments.Valid {
			err = json.Unmarshal([]byte(attachments.String), &attachmentsMetadata)
			if err != nil {
				log.Printf("failed to unmarshal attachments metadata: %v in prescription: %d", err, prescription.ID)
			} else {
				for _, attachment := range attachmentsMetadata {
					attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.PRESCRIPTION_BUCKET_NAME[len("maja."):])
				}
			}
			prescription.Attachments = attachmentsMetadata
		}
		prescriptions = append(prescriptions, &prescription)
	}
	return prescriptions, nil
}

func (r *PrescriptionRepositoryPostgresDB) Count(queries *models.PrescriptionsQueryRequestParams) (int64, error) {
	q := `SELECT
			COUNT(p.id)
		FROM prescriptions p
		LEFT JOIN customers c ON p.customerId = c.id
		LEFT JOIN users u ON c.userId = u.id`
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

func (r *PrescriptionRepositoryPostgresDB) Create(payload *models.PrescriptionsCreateRequestBody) (*domain.Prescription, error) {
	// Current time
	currentTime := time.Now()

	// Insert the prescription
	var insertedId int
	err := r.PostgresDB.QueryRow(`
		INSERT INTO prescriptions (customerId, title, datetime, doctorFullName, start_date, end_date, status, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		payload.CustomerID,
		payload.Title,
		payload.DateTimeAsDate,
		payload.DoctorFullName,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Status,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the prescription
	prescriptions, err := r.Query(&models.PrescriptionsQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(prescriptions) == 0 {
		return nil, errors.New("no prescription found")
	}
	prescription := prescriptions[0]

	// Return the prescription
	return prescription, nil
}

func (r *PrescriptionRepositoryPostgresDB) Delete(payload *models.PrescriptionsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM prescriptions
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

func (r *PrescriptionRepositoryPostgresDB) Update(payload *models.PrescriptionsUpdateRequestBody, id int) (*domain.Prescription, error) {
	var prescription domain.Prescription

	// Current time
	currentTime := time.Now()

	// Update the prescription
	var updatedId int
	err := r.PostgresDB.QueryRow(`
		UPDATE prescriptions
		SET customerId = $1, title = $2, datetime = $3, doctorFullName = $4, start_date = $5, end_date = $6, status = $7, updated_at = $8
		WHERE id = $9
		RETURNING id
	`,
		payload.CustomerID,
		payload.Title,
		payload.DateTimeAsDate,
		payload.DoctorFullName,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Status,
		currentTime,
		id,
	).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the prescription
	prescriptions, err := r.Query(&models.PrescriptionsQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(prescriptions) == 0 {
		return nil, errors.New("no prescription found")
	}
	prescription = *prescriptions[0]

	// Return the prescription
	return &prescription, nil
}

func (r *PrescriptionRepositoryPostgresDB) UpdatePrescriptionAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Prescription, error) {
	var prescription domain.Prescription

	// Current time
	currentTime := time.Now()

	// Find the prescription by id
	results, err := r.Query(&models.PrescriptionsQueryRequestParams{
		ID:    int(id),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	foundPrescription := results[0]
	if foundPrescription == nil {
		return nil, errors.New("no results found")
	}

	// Marshal attachments into JSON format
	for _, attachment := range previousAttachments {
		attachments = append(attachments, &attachment)
	}
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)

	// Update the prescription
	err = r.PostgresDB.QueryRow(`
		UPDATE prescriptions
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundPrescription.ID,
	).Scan(&prescription.ID)
	if err != nil {
		return nil, err
	}

	// Get the prescription
	prescriptions, err := r.Query(&models.PrescriptionsQueryRequestParams{
		ID:    int(foundPrescription.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(prescriptions) == 0 {
		return nil, errors.New("no prescription found")
	}
	prescription = *prescriptions[0]

	// Return the prescription
	return &prescription, nil
}
