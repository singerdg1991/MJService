package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/servicetype/domain"
	"github.com/hoitek/Maja-Service/internal/servicetype/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type ServiceTypeRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewServiceTypeRepositoryPostgresDB(d *sql.DB) *ServiceTypeRepositoryPostgresDB {
	return &ServiceTypeRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.ServiceTypesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" rc.id = %d", queries.ID))
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" rc.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.ServiceID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ServiceID.Op, fmt.Sprintf("%v", queries.Filters.ServiceID.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" rc.serviceID %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" rc.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *ServiceTypeRepositoryPostgresDB) Query(queries *models.ServiceTypesQueryRequestParams) ([]*domain.ServiceType, error) {
	q := `SELECT * FROM serviceTypes rc `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" rc.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" rc.description %s", queries.Sorts.Description.Op))
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

	var serviceTypes []*domain.ServiceType
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			serviceType domain.ServiceType
			service     domain.ServiceTypeService
		)
		err := rows.Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description, &serviceType.CreatedAt, &serviceType.UpdatedAt, &serviceType.DeletedAt)
		if err != nil {
			return nil, err
		}
		serviceRow := r.PostgresDB.QueryRow(`SELECT id, name, description FROM services WHERE id = $1`, serviceType.ServiceID)
		err = serviceRow.Scan(&service.ID, &service.Name, &service.Description)
		if err == nil {
			serviceType.Service = &domain.ServiceTypeService{
				ID:          service.ID,
				Name:        service.Name,
				Description: service.Description,
			}
		} else {
			log.Println(err)
		}

		serviceTypes = append(serviceTypes, &serviceType)
	}
	return serviceTypes, nil
}

func (r *ServiceTypeRepositoryPostgresDB) Count(queries *models.ServiceTypesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(rc.id) FROM serviceTypes rc `
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

func (r *ServiceTypeRepositoryPostgresDB) Create(payload *models.ServiceTypesCreateRequestBody) (*domain.ServiceType, error) {
	var serviceType domain.ServiceType

	// Current time
	currentTime := time.Now()

	// Insert the serviceType
	err := r.PostgresDB.QueryRow(`
		INSERT INTO serviceTypes (serviceId, name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, serviceId, name, description, created_at, updated_at, deleted_at
	`, payload.ServiceID, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description, &serviceType.CreatedAt, &serviceType.UpdatedAt, &serviceType.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Get the serviceType
	serviceTypes, err := r.Query(&models.ServiceTypesQueryRequestParams{
		ID: int(serviceType.ID),
	})
	if err != nil {
		return nil, err
	}
	if len(serviceTypes) == 0 {
		return nil, errors.New("no serviceType found")
	}
	serviceType = *serviceTypes[0]

	// Return the serviceType
	return &serviceType, nil
}

func (r *ServiceTypeRepositoryPostgresDB) Delete(payload *models.ServiceTypesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM serviceTypes
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

func (r *ServiceTypeRepositoryPostgresDB) Update(payload *models.ServiceTypesCreateRequestBody, id int64) (*domain.ServiceType, error) {
	var serviceType domain.ServiceType

	// Current time
	currentTime := time.Now()

	// Find the serviceType by id
	var foundServiceType domain.ServiceType
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM serviceTypes
		WHERE id = $1
	`, id).Scan(&foundServiceType.ID, &foundServiceType.ServiceID, &foundServiceType.Name, &foundServiceType.Description, &foundServiceType.CreatedAt, &foundServiceType.UpdatedAt, &foundServiceType.DeletedAt)

	// If the serviceType is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the serviceType
	err = r.PostgresDB.QueryRow(`
		UPDATE serviceTypes
		SET serviceId = $1, name = $2, updated_at = $3, description = $4
		WHERE id = $5
		RETURNING id, serviceId, name, description, created_at, updated_at, deleted_at
	`, payload.ServiceID, payload.Name, currentTime, payload.Description, foundServiceType.ID).Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description, &serviceType.CreatedAt, &serviceType.UpdatedAt, &serviceType.DeletedAt)

	// If the serviceType does not update, return an error
	if err != nil {
		return nil, err
	}

	// Get the serviceType
	serviceTypes, err := r.Query(&models.ServiceTypesQueryRequestParams{
		ID: int(serviceType.ID),
	})
	if err != nil {
		return nil, err
	}
	if len(serviceTypes) == 0 {
		return nil, errors.New("no serviceType found")
	}
	serviceType = *serviceTypes[0]

	// Return the serviceType
	return &serviceType, nil
}

func (r *ServiceTypeRepositoryPostgresDB) GetServiceTypesByIds(ids []int64) ([]*domain.ServiceType, error) {
	var serviceTypes []*domain.ServiceType
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM serviceTypes
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			serviceType domain.ServiceType
			service     domain.ServiceTypeService
		)
		err := rows.Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description, &serviceType.CreatedAt, &serviceType.UpdatedAt, &serviceType.DeletedAt)
		if err != nil {
			return nil, err
		}
		serviceRow := r.PostgresDB.QueryRow(`SELECT id, name, description FROM services WHERE id = $1`, serviceType.ServiceID)
		err = serviceRow.Scan(&service.ID, &service.Name, &service.Description)
		if err == nil {
			serviceType.Service = &domain.ServiceTypeService{
				ID:          service.ID,
				Name:        service.Name,
				Description: service.Description,
			}
		} else {
			log.Println(err)
		}
		serviceTypes = append(serviceTypes, &serviceType)
	}
	return serviceTypes, nil
}
