package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/trash/domain"
	"github.com/hoitek/Maja-Service/internal/trash/models"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type TrashRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewTrashRepositoryPostgresDB(d *sql.DB) *TrashRepositoryPostgresDB {
	return &TrashRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *TrashRepositoryPostgresDB) Query(queries *models.TrashesQueryRequestParams) ([]*domain.Trash, error) {
	q := `SELECT * FROM trashes `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit

		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return nil, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var trashes []*domain.Trash
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			trash       domain.Trash
			createdById sql.NullInt64
		)
		err := rows.Scan(&trash.ID, &trash.ModelName, &trash.ModelID, &trash.Reason, &trash.CreatedAt, &createdById)
		if err != nil {
			return nil, err
		}
		if createdById.Valid {
			userEntity, err := r.PostgresDB.Query("SELECT id, firstName, lastName FROM users WHERE id = $1", createdById.Int64)
			if err != nil {
				log.Println(err)
				userEntity.Close()
				continue
			}
			for userEntity.Next() {
				var (
					createdBy domain.TrashCreatedBy
				)
				err := userEntity.Scan(&createdBy.ID, &createdBy.FirstName, &createdBy.LastName)
				if err != nil {
					log.Println(err)
					userEntity.Close()
					continue
				}
				trash.CreatedBy = createdBy
			}
			userEntity.Close()
		}
		trashes = append(trashes, &trash)
	}
	return trashes, nil
}

func (r *TrashRepositoryPostgresDB) Count(queries *models.TrashesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM trashes `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit
		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return 0, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TrashRepositoryPostgresDB) Create(payload *models.TrashesCreateRequestBody) (*domain.Trash, error) {
	var trash domain.Trash

	// Current time
	currentTime := time.Now()

	// Insert the trash
	err := r.PostgresDB.QueryRow(`
		INSERT INTO trashes (modelName, modelId, reason, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, modelName, modelId, reason, created_at
	`, payload.ModelName, payload.ModelID, payload.Reason, currentTime, payload.CreatedBy.ID).Scan(&trash.ID, &trash.ModelName, &trash.ModelID, &trash.Reason, &trash.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Set the created by
	trash.CreatedBy = domain.TrashCreatedBy{
		ID:        payload.CreatedBy.ID,
		FirstName: payload.CreatedBy.FirstName,
		LastName:  payload.CreatedBy.LastName,
	}

	// Return the trash
	return &trash, nil
}

func (r *TrashRepositoryPostgresDB) Delete(payload *models.TrashesDeleteRequestBody) ([]int64, error) {
	var ids []int64
	idsStr := strings.Split(payload.IDs, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int64(id))
	}
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM trashes
		WHERE id = ANY ($1)
		RETURNING id
	`, pq.Int64Array(ids)).Scan(&rowsAffected)
	if err != nil {
		return nil, err
	}
	log.Println("rowsAffected", rowsAffected)
	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	return ids, nil
}
