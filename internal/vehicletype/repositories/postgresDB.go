package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/vehicletype/domain"
	"github.com/hoitek/Maja-Service/internal/vehicletype/models"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type VehicleTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewVehicleTypeRepositoryPostgresDB(d *sql.DB) *VehicleTypeRepositoryPostgresDB {
	return &VehicleTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *VehicleTypeRepositoryPostgresDB) Query(queries *models.VehicleTypesQueryRequestParams) ([]*domain.VehicleType, error) {
	q := `SELECT * FROM vehicletypes `
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

	var vehicletypes []*domain.VehicleType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var vehicletype domain.VehicleType
		err := rows.Scan(&vehicletype.ID, &vehicletype.Name, &vehicletype.CreatedAt, &vehicletype.UpdatedAt, &vehicletype.DeletedAt)
		if err != nil {
			return nil, err
		}
		vehicletypes = append(vehicletypes, &vehicletype)
	}
	return vehicletypes, nil
}

func (r *VehicleTypeRepositoryPostgresDB) Count(queries *models.VehicleTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM vehicletypes `
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

func (r *VehicleTypeRepositoryPostgresDB) Create(payload *models.VehicleTypesCreateRequestBody) (*domain.VehicleType, error) {
	var vehicletype domain.VehicleType

	// Current time
	currentTime := time.Now()

	// Find the vehicletype by name
	var foundVehicleType domain.VehicleType
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM vehicletypes
WHERE name = $1
`, payload.Name).Scan(&foundVehicleType.ID, &foundVehicleType.Name, &foundVehicleType.CreatedAt, &foundVehicleType.UpdatedAt, &foundVehicleType.DeletedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the vehicletype
	err = r.PostgresDB.QueryRow(`
INSERT INTO vehicletypes (name, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4)
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, currentTime, nil).Scan(&vehicletype.ID, &vehicletype.Name, &vehicletype.CreatedAt, &vehicletype.UpdatedAt, &vehicletype.DeletedAt)
	// If the vehicletype does not insert, return an error
	if err != nil {
		return nil, err
	}

	// Return the vehicletype
	return &vehicletype, nil
}

func (r *VehicleTypeRepositoryPostgresDB) Delete(payload *models.VehicleTypesDeleteRequestBody) ([]int64, error) {
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
DELETE FROM vehicletypes
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

func (r *VehicleTypeRepositoryPostgresDB) Update(payload *models.VehicleTypesCreateRequestBody, name string) (*domain.VehicleType, error) {
	var vehicletype domain.VehicleType

	// Current time
	currentTime := time.Now()

	// Find the vehicletype by name
	var foundVehicleType domain.VehicleType
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM vehicletypes
WHERE name = $1
`, name).Scan(&foundVehicleType.ID, &foundVehicleType.Name, &foundVehicleType.CreatedAt, &foundVehicleType.UpdatedAt, &foundVehicleType.DeletedAt)

	// If the vehicletype is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the vehicletype
	err = r.PostgresDB.QueryRow(`
UPDATE vehicletypes
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, foundVehicleType.ID).Scan(&vehicletype.ID, &vehicletype.Name, &vehicletype.CreatedAt, &vehicletype.UpdatedAt, &vehicletype.DeletedAt)

	// If the vehicletype does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the vehicletype
	return &vehicletype, nil
}
