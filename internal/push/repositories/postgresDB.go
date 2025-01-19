package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/push/domain"
	"github.com/hoitek/Maja-Service/internal/push/models"
	"github.com/lib/pq"
)

type PushRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewPushRepositoryPostgresDB(d *sql.DB) *PushRepositoryPostgresDB {
	return &PushRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.PushesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ls.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" ls.userId = %d", queries.UserID))
		}
	}
	return where
}

func (r *PushRepositoryPostgresDB) Query(queries *models.PushesQueryRequestParams) ([]*domain.Push, error) {
	q := `SELECT * FROM pushes ls `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var pushes []*domain.Push
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var push domain.Push
		err := rows.Scan(&push.ID, &push.UserID, &push.Endpoint, &push.KeysAuth, &push.KeysP256dh, &push.CreatedAt, &push.UpdatedAt, &push.DeletedAt)
		if err != nil {
			return nil, err
		}
		pushes = append(pushes, &push)
	}
	return pushes, nil
}

func (r *PushRepositoryPostgresDB) Count(queries *models.PushesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM pushes ls `
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

func (r *PushRepositoryPostgresDB) Create(payload *models.PushesCreateRequestBody) (*domain.Push, error) {
	var push domain.Push

	// Current time
	currentTime := time.Now()

	//  Check if the push already exists for the user
	var existingPushID int64
	err := r.PostgresDB.QueryRow(`
		SELECT id
		FROM pushes
		WHERE userId = $1
	`, payload.UserID).Scan(&existingPushID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingPushID != 0 {
		// Update the push
		err := r.PostgresDB.QueryRow(`
			UPDATE pushes
			SET endpoint = $1, keysAuth = $2, keysP256dh = $3, updated_at = $4
			WHERE id = $5
			RETURNING id, userId, endpoint, keysAuth, keysP256dh, created_at, updated_at, deleted_at
		`, payload.Endpoint, payload.KeysAuth, payload.KeysP256dh, currentTime, existingPushID).Scan(
			&push.ID,
			&push.UserID,
			&push.Endpoint,
			&push.KeysAuth,
			&push.KeysP256dh,
			&push.CreatedAt,
			&push.UpdatedAt,
			&push.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Insert the push
		err := r.PostgresDB.QueryRow(`
			INSERT INTO pushes (userId, endpoint, keysAuth, keysP256dh, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, userId, endpoint, keysAuth, keysP256dh, created_at, updated_at, deleted_at
		`, payload.UserID, payload.Endpoint, payload.KeysAuth, payload.KeysP256dh, currentTime, currentTime, nil).Scan(
			&push.ID,
			&push.UserID,
			&push.Endpoint,
			&push.KeysAuth,
			&push.KeysP256dh,
			&push.CreatedAt,
			&push.UpdatedAt,
			&push.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// Return the push
	return &push, nil
}

func (r *PushRepositoryPostgresDB) GetPushesByIds(ids []int64) ([]*domain.Push, error) {
	var pushes []*domain.Push
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM pushes
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var push domain.Push
		err := rows.Scan(&push.ID, &push.UserID, &push.Endpoint, &push.KeysAuth, &push.KeysP256dh, &push.CreatedAt, &push.UpdatedAt, &push.DeletedAt)
		if err != nil {
			return nil, err
		}
		pushes = append(pushes, &push)
	}
	return pushes, nil
}
