package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/limitation/domain"
	"github.com/hoitek/Maja-Service/internal/limitation/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type LimitationRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewLimitationRepositoryPostgresDB(d *sql.DB) *LimitationRepositoryPostgresDB {
	return &LimitationRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.LimitationsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ls.id = %d", queries.ID))
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *LimitationRepositoryPostgresDB) Query(queries *models.LimitationsQueryRequestParams) ([]*domain.Limitation, error) {
	q := `SELECT * FROM limitations ls `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.description %s", queries.Sorts.Description.Op))
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

	var limitations []*domain.Limitation
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var limitation domain.Limitation
		err := rows.Scan(&limitation.ID, &limitation.Name, &limitation.Description, &limitation.CreatedAt, &limitation.UpdatedAt, &limitation.DeletedAt)
		if err != nil {
			return nil, err
		}
		limitations = append(limitations, &limitation)
	}
	return limitations, nil
}

func (r *LimitationRepositoryPostgresDB) Count(queries *models.LimitationsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM limitations ls `
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

func (r *LimitationRepositoryPostgresDB) Create(payload *models.LimitationsCreateRequestBody) (*domain.Limitation, error) {
	var limitation domain.Limitation

	// Current time
	currentTime := time.Now()

	// Insert the limitation
	err := r.PostgresDB.QueryRow(`
		INSERT INTO limitations (name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&limitation.ID, &limitation.Name, &limitation.Description, &limitation.CreatedAt, &limitation.UpdatedAt, &limitation.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the limitation
	return &limitation, nil
}

func (r *LimitationRepositoryPostgresDB) Delete(payload *models.LimitationsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM limitations
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

func (r *LimitationRepositoryPostgresDB) Update(payload *models.LimitationsCreateRequestBody, id int64) (*domain.Limitation, error) {
	var limitation domain.Limitation

	// Current time
	currentTime := time.Now()

	// Find the limitation by name
	var foundLimitation domain.Limitation
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM limitations
		WHERE id = $1
	`, id).Scan(&foundLimitation.ID, &foundLimitation.Name, &foundLimitation.Description, &foundLimitation.CreatedAt, &foundLimitation.UpdatedAt, &foundLimitation.DeletedAt)

	// If the limitation is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the limitation
	err = r.PostgresDB.QueryRow(`
		UPDATE limitations
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, foundLimitation.ID).Scan(&limitation.ID, &limitation.Name, &limitation.Description, &limitation.CreatedAt, &limitation.UpdatedAt, &limitation.DeletedAt)

	// If the limitation does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the limitation
	return &limitation, nil
}

func (r *LimitationRepositoryPostgresDB) GetLimitationsByIds(ids []int64) ([]*domain.Limitation, error) {
	var limitations []*domain.Limitation
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM limitations
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var limitation domain.Limitation
		err := rows.Scan(&limitation.ID, &limitation.Name, &limitation.Description, &limitation.CreatedAt, &limitation.UpdatedAt, &limitation.DeletedAt)
		if err != nil {
			return nil, err
		}
		limitations = append(limitations, &limitation)
	}
	return limitations, nil
}
