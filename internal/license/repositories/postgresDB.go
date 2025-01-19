package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/license/domain"
	"github.com/hoitek/Maja-Service/internal/license/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type LicenseRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewLicenseRepositoryPostgresDB(d *sql.DB) *LicenseRepositoryPostgresDB {
	return &LicenseRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.LicensesQueryRequestParams) []string {
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

func (r *LicenseRepositoryPostgresDB) Query(queries *models.LicensesQueryRequestParams) ([]*domain.License, error) {
	q := `SELECT * FROM licenses ls `
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

	var licenses []*domain.License
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var license domain.License
		err := rows.Scan(&license.ID, &license.Name, &license.Description, &license.CreatedAt, &license.UpdatedAt, &license.DeletedAt)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, &license)
	}
	return licenses, nil
}

func (r *LicenseRepositoryPostgresDB) Count(queries *models.LicensesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM licenses ls `
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

func (r *LicenseRepositoryPostgresDB) Create(payload *models.LicensesCreateRequestBody) (*domain.License, error) {
	var license domain.License

	// Current time
	currentTime := time.Now()

	// Insert the license
	err := r.PostgresDB.QueryRow(`
		INSERT INTO licenses (name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&license.ID, &license.Name, &license.Description, &license.CreatedAt, &license.UpdatedAt, &license.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the license
	return &license, nil
}

func (r *LicenseRepositoryPostgresDB) Delete(payload *models.LicensesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM licenses
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

func (r *LicenseRepositoryPostgresDB) Update(payload *models.LicensesCreateRequestBody, id int64) (*domain.License, error) {
	var license domain.License

	// Current time
	currentTime := time.Now()

	// Find the license by name
	var foundLicense domain.License
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM licenses
		WHERE id = $1
	`, id).Scan(&foundLicense.ID, &foundLicense.Name, &foundLicense.Description, &foundLicense.CreatedAt, &foundLicense.UpdatedAt, &foundLicense.DeletedAt)

	// If the license is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the license
	err = r.PostgresDB.QueryRow(`
		UPDATE licenses
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, foundLicense.ID).Scan(&license.ID, &license.Name, &license.Description, &license.CreatedAt, &license.UpdatedAt, &license.DeletedAt)

	// If the license does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the license
	return &license, nil
}

func (r *LicenseRepositoryPostgresDB) GetLicensesByIds(ids []int64) ([]*domain.License, error) {
	var licenses []*domain.License
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM licenses
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var license domain.License
		err := rows.Scan(&license.ID, &license.Name, &license.Description, &license.CreatedAt, &license.UpdatedAt, &license.DeletedAt)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, &license)
	}
	return licenses, nil
}
