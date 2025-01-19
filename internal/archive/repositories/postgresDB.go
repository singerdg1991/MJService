package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/archive/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/archive/domain"
	"github.com/hoitek/Maja-Service/internal/archive/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type ArchiveRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewArchiveRepositoryPostgresDB(d *sql.DB) *ArchiveRepositoryPostgresDB {
	return &ArchiveRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.ArchivesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ar.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ar.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ar.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *ArchiveRepositoryPostgresDB) Query(queries *models.ArchivesQueryRequestParams) ([]*domain.Archive, error) {
	q := `
		SELECT
			ar.*,
			u.firstName as "userFirstName",
			u.lastName as "userLastName",
			u.avatarUrl as "userAvatarUrl"
		FROM
			archives ar
		LEFT JOIN users u ON u.id = ar.userId
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ar.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ar.description %s", queries.Sorts.Description.Op))
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

	var archives []*domain.Archive
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			archive             domain.Archive
			description         sql.NullString
			deletedAt           pq.NullTime
			userFirstName       sql.NullString
			userLastName        sql.NullString
			userAvatarUrl       sql.NullString
			attachments         json.RawMessage
			attachmentsMetadata []*types.UploadMetadata
		)
		err := rows.Scan(
			&archive.ID,
			&archive.UserID,
			&archive.Title,
			&archive.Subject,
			&description,
			&attachments,
			&archive.DateTime,
			&archive.CreatedAt,
			&archive.UpdatedAt,
			&deletedAt,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			archive.Description = &description.String
		}
		if deletedAt.Valid {
			archive.DeletedAt = &deletedAt.Time
		}
		archive.User = &domain.ArchiveUser{
			ID: archive.UserID,
		}
		if userFirstName.Valid {
			archive.User.FirstName = userFirstName.String
		}
		if userLastName.Valid {
			archive.User.LastName = userLastName.String
		}
		if userAvatarUrl.Valid {
			archive.User.AvatarUrl = userAvatarUrl.String
		}
		err = json.Unmarshal(attachments, &attachmentsMetadata)
		if err != nil {
			log.Printf("failed to unmarshal attachments metadata: %v in archive: %d", err, archive.ID)
		} else {
			for _, attachment := range attachmentsMetadata {
				attachment.Path = fmt.Sprintf("/%s/%s", "uploads", constants.ARCHIVE_BUCKET_NAME[len("maja."):])
			}
		}
		archive.Attachments = attachmentsMetadata
		archive.Date = archive.DateTime.Format("2006-01-02")
		archive.Time = archive.DateTime.Format("15:04")
		archives = append(archives, &archive)
	}
	return archives, nil
}

func (r *ArchiveRepositoryPostgresDB) Count(queries *models.ArchivesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ar.id) FROM archives ar `
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

func (r *ArchiveRepositoryPostgresDB) Create(payload *models.ArchivesCreateRequestBody) (*domain.Archive, error) {
	var archive domain.Archive

	// Current time
	currentTime := time.Now()

	// Insert the archive
	err := r.PostgresDB.QueryRow(`
		INSERT INTO archives (userId, title, subject, description, datetime, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.UserID,
		payload.Title,
		payload.Subject,
		payload.Description,
		payload.DateTime,
		currentTime,
		currentTime,
		nil,
	).Scan(
		&archive.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the archive
	results, err := r.Query(&models.ArchivesQueryRequestParams{
		ID:    int(archive.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	archive = *results[0]

	// Return the archive
	return &archive, nil
}

func (r *ArchiveRepositoryPostgresDB) Delete(payload *models.ArchivesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM archives
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

func (r *ArchiveRepositoryPostgresDB) Update(payload *models.ArchivesCreateRequestBody, id int64) (*domain.Archive, error) {
	var archive domain.Archive

	// Current time
	currentTime := time.Now()

	// Find the archive by id
	results, err := r.Query(&models.ArchivesQueryRequestParams{
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
	foundArchive := results[0]
	if foundArchive == nil {
		return nil, errors.New("no results found")
	}

	// Update the archive
	err = r.PostgresDB.QueryRow(`
		UPDATE archives
		SET userId = $1, title = $2, subject = $3, description = $4, datetime = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, userId, title, subject, description, datetime, created_at, updated_at, deleted_at
	`,
		payload.UserID,
		payload.Title,
		payload.Subject,
		payload.Description,
		payload.DateTime,
		currentTime,
		foundArchive.ID,
	).Scan(
		&archive.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the archive
	results, err = r.Query(&models.ArchivesQueryRequestParams{
		ID:    int(foundArchive.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	archive = *results[0]

	// Return the archive
	return &archive, nil
}

func (r *ArchiveRepositoryPostgresDB) GetArchivesByIds(ids []int64) ([]*domain.Archive, error) {
	var archives []*domain.Archive
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM archives ar
		LEFT JOIN users u ON u.id = ar.userIdq
		WHERE ar.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			archive       domain.Archive
			description   sql.NullString
			deletedAt     pq.NullTime
			userFirstName sql.NullString
			userLastName  sql.NullString
			userAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&archive.ID,
			&archive.UserID,
			&archive.Title,
			&archive.Subject,
			&description,
			&archive.DateTime,
			&archive.CreatedAt,
			&archive.UpdatedAt,
			&deletedAt,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			archive.Description = &description.String
		}
		if deletedAt.Valid {
			archive.DeletedAt = &deletedAt.Time
		}
		archive.User = &domain.ArchiveUser{
			ID: archive.UserID,
		}
		if userFirstName.Valid {
			archive.User.FirstName = userFirstName.String
		}
		if userLastName.Valid {
			archive.User.LastName = userLastName.String
		}
		if userAvatarUrl.Valid {
			archive.User.AvatarUrl = userAvatarUrl.String
		}
		archives = append(archives, &archive)
	}
	return archives, nil
}

func (r *ArchiveRepositoryPostgresDB) UpdateAttachments(attachments []*types.UploadMetadata, id int64) (*domain.Archive, error) {
	var archive domain.Archive

	// Current time
	currentTime := time.Now()

	// Find the archive by id
	results, err := r.Query(&models.ArchivesQueryRequestParams{
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
	foundArchive := results[0]
	if foundArchive == nil {
		return nil, errors.New("no results found")
	}

	// Marshal attachments into JSON format
	b, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	attachmentsJSON := string(b)

	// Update the archive
	err = r.PostgresDB.QueryRow(`
		UPDATE archives
		SET attachments = $1, updated_at = $2
		WHERE id = $3
		RETURNING id
	`,
		attachmentsJSON,
		currentTime,
		foundArchive.ID,
	).Scan(
		&archive.ID,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the archive
	results, err = r.Query(&models.ArchivesQueryRequestParams{
		ID:    int(foundArchive.ID),
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	archive = *results[0]

	// Return the archive
	return &archive, nil
}
