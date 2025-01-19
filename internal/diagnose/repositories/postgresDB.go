package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/diagnose/domain"
	"github.com/hoitek/Maja-Service/internal/diagnose/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type DiagnoseRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewDiagnoseRepositoryPostgresDB(d *sql.DB) *DiagnoseRepositoryPostgresDB {
	return &DiagnoseRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.DiagnosesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" d.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" d.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" d.description %s %s", opValue.Operator, val))
		}
		if queries.Filters.Code.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Code.Op, fmt.Sprintf("%v", queries.Filters.Code.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" d.code %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" d.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *DiagnoseRepositoryPostgresDB) Query(queries *models.DiagnosesQueryRequestParams) ([]*domain.Diagnose, error) {
	q := `SELECT d.id, d.title, d.code, d.description, d.created_at, d.updated_at, d.deleted_at FROM diagnoses d`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" d.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.Code.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" d.code %s", queries.Sorts.Code.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" d.description %s", queries.Sorts.Description.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" d.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var diagnoses []*domain.Diagnose
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			diagnose    domain.Diagnose
			description sql.NullString
			deleteAt    pq.NullTime
		)
		err := rows.Scan(
			&diagnose.ID,
			&diagnose.Title,
			&diagnose.Code,
			&description,
			&diagnose.CreatedAt,
			&diagnose.UpdatedAt,
			&deleteAt,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			diagnose.Description = &description.String
		}
		if deleteAt.Valid {
			diagnose.DeletedAt = &deleteAt.Time
		}
		diagnoses = append(diagnoses, &diagnose)
	}
	return diagnoses, nil
}

func (r *DiagnoseRepositoryPostgresDB) Count(queries *models.DiagnosesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(d.id) FROM diagnoses d `
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

func (r *DiagnoseRepositoryPostgresDB) Create(payload *models.DiagnosesCreateRequestBody) (*domain.Diagnose, error) {
	var diagnose domain.Diagnose

	// Current time
	currentTime := time.Now()

	// Insert the diagnose
	var (
		description sql.NullString
		deletedAt   pq.NullTime
	)
	err := r.PostgresDB.QueryRow(`
		INSERT INTO diagnoses (title, code, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, code, description, created_at, updated_at, deleted_at
	`, payload.Title, payload.Code, payload.Description, currentTime, currentTime, nil).Scan(&diagnose.ID, &diagnose.Title, &diagnose.Code, &description, &diagnose.CreatedAt, &diagnose.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		diagnose.Description = &description.String
	}
	if deletedAt.Valid {
		diagnose.DeletedAt = &deletedAt.Time
	}

	// Return the diagnose
	return &diagnose, nil
}

func (r *DiagnoseRepositoryPostgresDB) Delete(payload *models.DiagnosesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM diagnoses
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

func (r *DiagnoseRepositoryPostgresDB) Update(payload *models.DiagnosesCreateRequestBody, id int) (*domain.Diagnose, error) {
	var diagnose domain.Diagnose

	// Current time
	currentTime := time.Now()

	// Find the diagnose by id
	var (
		foundDiagnose domain.Diagnose
		description   sql.NullString
		deletedAt     pq.NullTime
	)
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM diagnoses
		WHERE id = $1
	`, id).Scan(&foundDiagnose.ID, &foundDiagnose.Title, &foundDiagnose.Code, &description, &foundDiagnose.CreatedAt, &foundDiagnose.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		foundDiagnose.Description = &description.String
	}
	if deletedAt.Valid {
		foundDiagnose.DeletedAt = &deletedAt.Time
	}

	// Update the diagnose
	var (
		updatedDescription sql.NullString
		updateDeletedAt    pq.NullTime
	)
	err = r.PostgresDB.QueryRow(`
		UPDATE diagnoses
		SET title = $1, code = $2, description = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, code, description, created_at, updated_at, deleted_at
	`, payload.Title, payload.Code, payload.Description, currentTime, foundDiagnose.ID).Scan(
		&diagnose.ID,
		&diagnose.Title,
		&diagnose.Code,
		&updatedDescription,
		&diagnose.CreatedAt,
		&diagnose.UpdatedAt,
		&updateDeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if updatedDescription.Valid {
		diagnose.Description = &updatedDescription.String
	}
	if updateDeletedAt.Valid {
		diagnose.DeletedAt = &updateDeletedAt.Time
	}

	// Return the diagnose
	return &diagnose, nil
}
