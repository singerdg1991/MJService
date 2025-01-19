package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/servicegrade/domain"
	"github.com/hoitek/Maja-Service/internal/servicegrade/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type ServiceGradeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewServiceGradeRepositoryPostgresDB(d *sql.DB) *ServiceGradeRepositoryPostgresDB {
	return &ServiceGradeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.ServiceGradesQueryRequestParams) []string {
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
		if queries.Filters.Grade.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Grade.Op, fmt.Sprintf("%v", queries.Filters.Grade.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.grade %s %s", opValue.Operator, val))
		}
		if queries.Filters.Color.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Color.Op, fmt.Sprintf("%v", queries.Filters.Color.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.color %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *ServiceGradeRepositoryPostgresDB) Query(queries *models.ServiceGradesQueryRequestParams) ([]*domain.ServiceGrade, error) {
	q := `SELECT * FROM servicegrades ls `
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
		if queries.Sorts.Grade.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.grade %s", queries.Sorts.Grade.Op))
		}
		if queries.Sorts.Color.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.color %s", queries.Sorts.Color.Op))
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

	var servicegrades []*domain.ServiceGrade
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var servicegrade domain.ServiceGrade
		err := rows.Scan(&servicegrade.ID, &servicegrade.Name, &servicegrade.Description, &servicegrade.Grade, &servicegrade.Color, &servicegrade.CreatedAt, &servicegrade.UpdatedAt, &servicegrade.DeletedAt)
		if err != nil {
			return nil, err
		}
		servicegrades = append(servicegrades, &servicegrade)
	}
	return servicegrades, nil
}

func (r *ServiceGradeRepositoryPostgresDB) Count(queries *models.ServiceGradesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM servicegrades ls `
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

func (r *ServiceGradeRepositoryPostgresDB) Create(payload *models.ServiceGradesCreateRequestBody) (*domain.ServiceGrade, error) {
	var servicegrade domain.ServiceGrade

	// Current time
	currentTime := time.Now()

	// Insert the servicegrade
	err := r.PostgresDB.QueryRow(`
		INSERT INTO servicegrades (name, description, grade, color, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, grade, color, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, payload.Grade, payload.Color, currentTime, currentTime, nil).Scan(&servicegrade.ID, &servicegrade.Name, &servicegrade.Description, &servicegrade.Grade, &servicegrade.Color, &servicegrade.CreatedAt, &servicegrade.UpdatedAt, &servicegrade.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the servicegrade
	return &servicegrade, nil
}

func (r *ServiceGradeRepositoryPostgresDB) Delete(payload *models.ServiceGradesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM servicegrades
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

func (r *ServiceGradeRepositoryPostgresDB) Update(payload *models.ServiceGradesCreateRequestBody, id int64) (*domain.ServiceGrade, error) {
	var servicegrade domain.ServiceGrade

	// Current time
	currentTime := time.Now()

	// Find the servicegrade by name
	var foundServiceGrade domain.ServiceGrade
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM servicegrades
		WHERE id = $1
	`, id).Scan(&foundServiceGrade.ID, &foundServiceGrade.Name, &foundServiceGrade.Description, &foundServiceGrade.Grade, &foundServiceGrade.Color, &foundServiceGrade.CreatedAt, &foundServiceGrade.UpdatedAt, &foundServiceGrade.DeletedAt)

	// If the servicegrade is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the servicegrade
	err = r.PostgresDB.QueryRow(`
		UPDATE servicegrades
		SET name = $1, updated_at = $2, description = $3, grade = $4, color = $5
		WHERE id = $6
		RETURNING id, name, description, grade, color, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, payload.Grade, payload.Color, foundServiceGrade.ID).Scan(&servicegrade.ID, &servicegrade.Name, &servicegrade.Description, &servicegrade.Grade, &servicegrade.Color, &servicegrade.CreatedAt, &servicegrade.UpdatedAt, &servicegrade.DeletedAt)

	// If the servicegrade does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the servicegrade
	return &servicegrade, nil
}

func (r *ServiceGradeRepositoryPostgresDB) GetServiceGradesByIds(ids []int64) ([]*domain.ServiceGrade, error) {
	var servicegrades []*domain.ServiceGrade
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM servicegrades
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var servicegrade domain.ServiceGrade
		err := rows.Scan(&servicegrade.ID, &servicegrade.Name, &servicegrade.Grade, &servicegrade.Description, &servicegrade.CreatedAt, &servicegrade.UpdatedAt, &servicegrade.DeletedAt)
		if err != nil {
			return nil, err
		}
		servicegrades = append(servicegrades, &servicegrade)
	}
	return servicegrades, nil
}
