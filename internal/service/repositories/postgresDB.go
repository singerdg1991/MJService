package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/service/domain"
	"github.com/hoitek/Maja-Service/internal/service/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type ServiceRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewServiceRepositoryPostgresDB(d *sql.DB) *ServiceRepositoryPostgresDB {
	return &ServiceRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.ServicesQueryRequestParams) []string {
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

func (r *ServiceRepositoryPostgresDB) Query(queries *models.ServicesQueryRequestParams) ([]*domain.Service, error) {
	q := `SELECT * FROM services ls `
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

	var services []*domain.Service
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var service domain.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt, &service.DeletedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	return services, nil
}

func (r *ServiceRepositoryPostgresDB) Count(queries *models.ServicesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM services ls `
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

func (r *ServiceRepositoryPostgresDB) Create(payload *models.ServicesCreateRequestBody) (*domain.Service, error) {
	var service domain.Service

	// Current time
	currentTime := time.Now()

	// Insert the service
	err := r.PostgresDB.QueryRow(`
		INSERT INTO services (name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt, &service.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the service
	return &service, nil
}

func (r *ServiceRepositoryPostgresDB) Delete(payload *models.ServicesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM services
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

func (r *ServiceRepositoryPostgresDB) Update(payload *models.ServicesCreateRequestBody, id int64) (*domain.Service, error) {
	var service domain.Service

	// Current time
	currentTime := time.Now()

	// Find the service by name
	var foundService domain.Service
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM services
		WHERE id = $1
	`, id).Scan(&foundService.ID, &foundService.Name, &foundService.Description, &foundService.CreatedAt, &foundService.UpdatedAt, &foundService.DeletedAt)

	// If the service is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the service
	err = r.PostgresDB.QueryRow(`
		UPDATE services
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, foundService.ID).Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt, &service.DeletedAt)

	// If the service does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the service
	return &service, nil
}

func (r *ServiceRepositoryPostgresDB) GetServicesByIds(ids []int64) ([]*domain.Service, error) {
	var services []*domain.Service
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM services
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var service domain.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt, &service.DeletedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	return services, nil
}
