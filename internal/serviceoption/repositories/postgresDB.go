package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/serviceoption/domain"
	"github.com/hoitek/Maja-Service/internal/serviceoption/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type ServiceOptionRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewServiceOptionRepositoryPostgresDB(d *sql.DB) *ServiceOptionRepositoryPostgresDB {
	return &ServiceOptionRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.ServiceOptionsQueryRequestParams) []string {
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
		if queries.Filters.ServiceTypeID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ServiceTypeID.Op, fmt.Sprintf("%v", queries.Filters.ServiceTypeID.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" rc.serviceTypeID %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" rc.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *ServiceOptionRepositoryPostgresDB) Query(queries *models.ServiceOptionsQueryRequestParams) ([]*domain.ServiceOption, error) {
	q := `SELECT * FROM serviceOptions rc `
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

	var serviceOptions []*domain.ServiceOption
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			serviceoption domain.ServiceOption
			serviceType   domain.ServiceOptionServiceType
			service       domain.ServiceOptionServiceTypeService
		)
		err := rows.Scan(&serviceoption.ID, &serviceoption.ServiceTypeID, &serviceoption.Name, &serviceoption.Description, &serviceoption.CreatedAt, &serviceoption.UpdatedAt, &serviceoption.DeletedAt)
		if err != nil {
			return nil, err
		}
		serviceTypeRow := r.PostgresDB.QueryRow(`SELECT id, serviceId, name, description FROM serviceTypes WHERE id = $1`, serviceoption.ServiceTypeID)
		err = serviceTypeRow.Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description)
		if err == nil {
			serviceoption.ServiceType = &domain.ServiceOptionServiceType{
				ID:          serviceType.ID,
				ServiceID:   serviceType.ServiceID,
				Name:        serviceType.Name,
				Description: serviceType.Description,
			}
			serviceRow := r.PostgresDB.QueryRow(`SELECT id, name, description FROM services WHERE id = $1`, serviceType.ServiceID)
			err = serviceRow.Scan(&service.ID, &service.Name, &service.Description)
			if err == nil {
				serviceoption.ServiceType.Service = domain.ServiceOptionServiceTypeService{
					ID:          service.ID,
					Name:        service.Name,
					Description: service.Description,
				}
			} else {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
		serviceOptions = append(serviceOptions, &serviceoption)
	}
	return serviceOptions, nil
}

func (r *ServiceOptionRepositoryPostgresDB) Count(queries *models.ServiceOptionsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(rc.id) FROM serviceOptions rc `
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

func (r *ServiceOptionRepositoryPostgresDB) Create(payload *models.ServiceOptionsCreateRequestBody) (*domain.ServiceOption, error) {
	var serviceType domain.ServiceOption

	// Current time
	currentTime := time.Now()

	// Insert the serviceType
	err := r.PostgresDB.QueryRow(`
		INSERT INTO serviceOptions (serviceTypeId, name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, serviceTypeId, name, description, created_at, updated_at, deleted_at
	`, payload.ServiceTypeID, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&serviceType.ID, &serviceType.ServiceTypeID, &serviceType.Name, &serviceType.Description, &serviceType.CreatedAt, &serviceType.UpdatedAt, &serviceType.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Get the serviceType
	services, err := r.Query(&models.ServiceOptionsQueryRequestParams{
		ID: int(serviceType.ID),
	})
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, errors.New("no serviceType found")
	}
	serviceType = *services[0]

	// Return the serviceType
	return &serviceType, nil
}

func (r *ServiceOptionRepositoryPostgresDB) Delete(payload *models.ServiceOptionsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM serviceOptions
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

func (r *ServiceOptionRepositoryPostgresDB) Update(payload *models.ServiceOptionsCreateRequestBody, id int64) (*domain.ServiceOption, error) {
	var serviceOption domain.ServiceOption

	// Current time
	currentTime := time.Now()

	// Find the serviceType by id
	var foundServiceOption domain.ServiceOption
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM serviceOptions
		WHERE id = $1
	`, id).Scan(&foundServiceOption.ID, &foundServiceOption.ServiceTypeID, &foundServiceOption.Name, &foundServiceOption.Description, &foundServiceOption.CreatedAt, &foundServiceOption.UpdatedAt, &foundServiceOption.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Update the serviceType
	err = r.PostgresDB.QueryRow(`
		UPDATE serviceOptions
		SET serviceTypeId = $1, name = $2, updated_at = $3, description = $4
		WHERE id = $5
		RETURNING id, serviceTypeId, name, description, created_at, updated_at, deleted_at
	`, payload.ServiceTypeID, payload.Name, currentTime, payload.Description, foundServiceOption.ID).Scan(&serviceOption.ID, &serviceOption.ServiceTypeID, &serviceOption.Name, &serviceOption.Description, &serviceOption.CreatedAt, &serviceOption.UpdatedAt, &serviceOption.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Get the serviceType
	serviceOptions, err := r.Query(&models.ServiceOptionsQueryRequestParams{
		ID: int(serviceOption.ID),
	})
	if err != nil {
		return nil, err
	}
	if len(serviceOptions) == 0 {
		return nil, errors.New("no serviceType found")
	}
	serviceOption = *serviceOptions[0]

	// Return the serviceOption
	return &serviceOption, nil
}

func (r *ServiceOptionRepositoryPostgresDB) GetServiceOptionsByIds(ids []int64) ([]*domain.ServiceOption, error) {
	var serviceOptions []*domain.ServiceOption
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM serviceOptions
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			serviceoption domain.ServiceOption
			serviceType   domain.ServiceOptionServiceType
			service       domain.ServiceOptionServiceTypeService
		)
		err := rows.Scan(&serviceoption.ID, &serviceoption.ServiceTypeID, &serviceoption.Name, &serviceoption.Description, &serviceoption.CreatedAt, &serviceoption.UpdatedAt, &serviceoption.DeletedAt)
		if err != nil {
			return nil, err
		}
		serviceTypeRow := r.PostgresDB.QueryRow(`SELECT id, serviceId, name, description FROM serviceTypes WHERE id = $1`, serviceoption.ServiceTypeID)
		err = serviceTypeRow.Scan(&serviceType.ID, &serviceType.ServiceID, &serviceType.Name, &serviceType.Description)
		if err == nil {
			serviceoption.ServiceType = &domain.ServiceOptionServiceType{
				ID:          serviceType.ID,
				ServiceID:   serviceType.ServiceID,
				Name:        serviceType.Name,
				Description: serviceType.Description,
			}
			serviceRow := r.PostgresDB.QueryRow(`SELECT id, name, description FROM services WHERE id = $1`, serviceType.ServiceID)
			err = serviceRow.Scan(&service.ID, &service.Name, &service.Description)
			if err == nil {
				serviceoption.ServiceType.Service = domain.ServiceOptionServiceTypeService{
					ID:          service.ID,
					Name:        service.Name,
					Description: service.Description,
				}
			} else {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
		serviceOptions = append(serviceOptions, &serviceoption)
	}
	return serviceOptions, nil
}
