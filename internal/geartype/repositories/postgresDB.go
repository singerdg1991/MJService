package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/geartype/domain"
	"github.com/hoitek/Maja-Service/internal/geartype/models"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type GearTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewGearTypeRepositoryPostgresDB(d *sql.DB) *GearTypeRepositoryPostgresDB {
	return &GearTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *GearTypeRepositoryPostgresDB) Query(queries *models.GearTypesQueryRequestParams) ([]*domain.GearType, error) {
	q := `SELECT * FROM geartypes `
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

	var geartypes []*domain.GearType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var geartype domain.GearType
		err := rows.Scan(&geartype.ID, &geartype.Name, &geartype.CreatedAt, &geartype.UpdatedAt, &geartype.DeletedAt)
		if err != nil {
			return nil, err
		}
		geartypes = append(geartypes, &geartype)
	}
	return geartypes, nil
}

func (r *GearTypeRepositoryPostgresDB) Count(queries *models.GearTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM geartypes `
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

func (r *GearTypeRepositoryPostgresDB) Create(payload *models.GearTypesCreateRequestBody) (*domain.GearType, error) {
	var geartype domain.GearType

	// Current time
	currentTime := time.Now()

	// Find the geartype by name
	var foundGearType domain.GearType
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM geartypes
WHERE name = $1
`, payload.Name).Scan(&foundGearType.ID, &foundGearType.Name, &foundGearType.CreatedAt, &foundGearType.UpdatedAt, &foundGearType.DeletedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the geartype
	err = r.PostgresDB.QueryRow(`
INSERT INTO geartypes (name, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4)
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, currentTime, nil).Scan(&geartype.ID, &geartype.Name, &geartype.CreatedAt, &geartype.UpdatedAt, &geartype.DeletedAt)
	// If the geartype does not insert, return an error
	if err != nil {
		return nil, err
	}

	// Return the geartype
	return &geartype, nil
}

func (r *GearTypeRepositoryPostgresDB) Delete(payload *models.GearTypesDeleteRequestBody) ([]int64, error) {
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
DELETE FROM geartypes
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

func (r *GearTypeRepositoryPostgresDB) Update(payload *models.GearTypesCreateRequestBody, name string) (*domain.GearType, error) {
	var geartype domain.GearType

	// Current time
	currentTime := time.Now()

	// Find the geartype by name
	var foundGearType domain.GearType
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM geartypes
WHERE name = $1
`, name).Scan(&foundGearType.ID, &foundGearType.Name, &foundGearType.CreatedAt, &foundGearType.UpdatedAt, &foundGearType.DeletedAt)

	// If the geartype is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the geartype
	err = r.PostgresDB.QueryRow(`
UPDATE geartypes
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, foundGearType.ID).Scan(&geartype.ID, &geartype.Name, &geartype.CreatedAt, &geartype.UpdatedAt, &geartype.DeletedAt)

	// If the geartype does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the geartype
	return &geartype, nil
}
