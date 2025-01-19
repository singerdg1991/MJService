package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/address/domain"
	"github.com/hoitek/Maja-Service/internal/address/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type AddressRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewAddressRepositoryPostgresDB(d *sql.DB) *AddressRepositoryPostgresDB {
	return &AddressRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.AddressesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf("a.id = %d", queries.ID))
		}
		if queries.StaffID != 0 {
			where = append(where, fmt.Sprintf("a.staffId = %d", queries.StaffID))
		}
		if queries.CustomerID != 0 {
			where = append(where, fmt.Sprintf("a.customerId = %d", queries.CustomerID))
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" a.name %s %s ", opValue.Operator, val))
		}
		if queries.Filters.Street.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Street.Op, fmt.Sprintf("%v", queries.Filters.Street.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" a.street %s %s ", opValue.Operator, val))
		}
		if queries.Filters.PostalCode.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PostalCode.Op, fmt.Sprintf("%v", queries.Filters.PostalCode.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" a.postalCode %s %s ", opValue.Operator, val))
		}
		if queries.Filters.BuildingNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.BuildingNumber.Op, fmt.Sprintf("%v", queries.Filters.BuildingNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" a.buildingNumber %s %s ", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" a.created_at %s %s ", opValue.Operator, val))
		}
	}
	return where
}

func (r *AddressRepositoryPostgresDB) Query(queries *models.AddressesQueryRequestParams) ([]*domain.Address, error) {
	q := `
		SELECT
			a.id AS id,
			a.staffId AS staffId,
			a.customerId AS customerId,
			a.cityId AS cityId,
			a.street AS street,
			a.name AS name,
			a.postalCode AS postalCode,
			a.buildingNumber AS buildingNumber,
			a.isDeliveryAddress AS isDeliveryAddress,
			a.isMainAddress AS isMainAddress,
			a.created_at AS created_at,
			a.updated_at AS updated_at,
			a.deleted_at AS deleted_at,
			u.firstName AS staffFirstName,
			u.lastName AS staffLastName,
			u.email AS staffEmail,
			u2.firstName AS customerFirstName,
			u2.lastName AS customerLastName,
			u2.email AS customerEmail
		FROM addresses a
		LEFT JOIN staffs s ON a.staffId = s.id
		LEFT JOIN customers c ON a.customerId = c.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN users u2 ON c.userId = u2.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}

		var sorts []string
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" a.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Street.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" a.street %s", queries.Sorts.Street.Op))
		}
		if queries.Sorts.PostalCode.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" a.postalCode %s", queries.Sorts.PostalCode.Op))
		}
		if queries.Sorts.BuildingNumber.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" a.buildingNumber %s", queries.Sorts.BuildingNumber.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" a.created_at %s", queries.Sorts.CreatedAt.Op))
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
	log.Println(q)

	var addresses []*domain.Address
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			address           domain.Address
			staff             *domain.AddressStaff
			customer          *domain.AddressCustomer
			staffId           sql.NullInt64
			staffFirstName    sql.NullString
			staffLastName     sql.NullString
			staffEmail        sql.NullString
			customerId        sql.NullInt64
			customerFirstName sql.NullString
			customerLastName  sql.NullString
			customerEmail     sql.NullString
			deletedAt         sql.NullTime
		)
		err := rows.Scan(
			&address.ID,
			&staffId,
			&customerId,
			&address.CityID,
			&address.Street,
			&address.Name,
			&address.PostalCode,
			&address.BuildingNumber,
			&address.IsDeliveryAddress,
			&address.IsMainAddress,
			&address.CreatedAt,
			&address.UpdatedAt,
			&deletedAt,
			&staffFirstName,
			&staffLastName,
			&staffEmail,
			&customerFirstName,
			&customerLastName,
			&customerEmail,
		)
		if err != nil {
			return nil, err
		}
		if staffId.Valid {
			sid := uint(staffId.Int64)
			staff = &domain.AddressStaff{
				ID:        sid,
				FirstName: staffFirstName.String,
				LastName:  staffLastName.String,
				Email:     staffEmail.String,
			}
			address.StaffID = &sid
		}
		if customerId.Valid {
			cid := uint(customerId.Int64)
			customer = &domain.AddressCustomer{
				ID:        cid,
				FirstName: customerFirstName.String,
				LastName:  customerLastName.String,
				Email:     customerEmail.String,
			}
			address.CustomerID = &cid
		}
		if deletedAt.Valid {
			address.DeletedAt = &deletedAt.Time
		}
		address.Staff = staff
		address.Customer = customer

		// Get the city
		address.City = &domain.AddressCity{}
		err = r.PostgresDB.QueryRow(`
			SELECT id, name
			FROM cities
			WHERE id = $1
		`, address.CityID).Scan(
			&address.City.ID,
			&address.City.Name,
		)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, &address)
	}
	return addresses, nil
}

func (r *AddressRepositoryPostgresDB) Count(queries *models.AddressesQueryRequestParams) (int64, error) {
	q := `
		SELECT
			COUNT(a.id)
		FROM addresses a
		LEFT JOIN staffs s ON a.staffId = s.id
		LEFT JOIN customers c ON a.customerId = c.id
		LEFT JOIN users u ON s.userId = u.id
		LEFT JOIN users u2 ON c.userId = u2.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
	}
	q += ";"

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AddressRepositoryPostgresDB) Create(payload *models.AddressesCreateRequestBody) (*domain.Address, error) {
	var address domain.Address

	// Current time
	currentTime := time.Now()

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Reset isDeliveryAddress based on staff id or customer id
	if payload.IsDeliveryAddressAsBool {
		if payload.StaffID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isDeliveryAddress = false
				WHERE staffId = $1
			`, payload.StaffID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if payload.CustomerID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isDeliveryAddress = false
				WHERE customerId = $1
			`, payload.CustomerID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Reset isMainAddress based on staff id or customer id
	if payload.IsMainAddressAsBool {
		if payload.StaffID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isMainAddress = false
				WHERE staffId = $1
			`, payload.StaffID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if payload.CustomerID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isMainAddress = false
				WHERE customerId = $1
			`, payload.CustomerID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Insert the address
	err = tx.QueryRow(`
		INSERT INTO addresses (staffId, customerId, cityId, street, name, postalCode, buildingNumber, isDeliveryAddress, isMainAddress, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, staffId, customerId, cityId, street, name, postalCode, buildingNumber, isDeliveryAddress, isMainAddress, created_at, updated_at, deleted_at
	`,
		payload.StaffID,
		payload.CustomerID,
		payload.City.ID,
		payload.Street,
		payload.Name,
		payload.PostalCode,
		payload.BuildingNumber,
		payload.IsDeliveryAddressAsBool,
		payload.IsMainAddressAsBool,
		currentTime,
		currentTime,
		nil,
	).Scan(
		&address.ID,
		&address.StaffID,
		&address.CustomerID,
		&address.CityID,
		&address.Street,
		&address.Name,
		&address.PostalCode,
		&address.BuildingNumber,
		&address.IsDeliveryAddress,
		&address.IsMainAddress,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Get the address
	addresses, err := r.Query(&models.AddressesQueryRequestParams{
		ID: int(address.ID),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(addresses) == 0 {
		tx.Rollback()
		return nil, errors.New("something went wrong retrieving the address")
	}
	address = *addresses[0]

	// Return the address
	return &address, nil
}

func (r *AddressRepositoryPostgresDB) Delete(payload *models.AddressesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM addresses
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

func (r *AddressRepositoryPostgresDB) Update(payload *models.AddressesCreateRequestBody, id int) (*domain.Address, error) {
	var address domain.Address

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Current time
	currentTime := time.Now()

	// Find the address by name
	var foundAddress domain.Address
	err = tx.QueryRow(`
		SELECT *
		FROM addresses
		WHERE id = $1
	`, id).Scan(&foundAddress.ID, &foundAddress.StaffID, &foundAddress.CustomerID, &foundAddress.CityID, &foundAddress.Street, &foundAddress.Name, &foundAddress.PostalCode, &foundAddress.BuildingNumber, &foundAddress.IsDeliveryAddress, &foundAddress.IsMainAddress, &foundAddress.CreatedAt, &foundAddress.UpdatedAt, &foundAddress.DeletedAt)

	// If the address is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Reset isDeliveryAddress based on staff id or customer id
	if payload.IsDeliveryAddressAsBool {
		if payload.StaffID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isDeliveryAddress = false
				WHERE staffId = $1
			`, payload.StaffID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if payload.CustomerID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isDeliveryAddress = false
				WHERE customerId = $1
			`, payload.CustomerID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Reset isMainAddress based on staff id or customer id
	if payload.IsMainAddressAsBool {
		if payload.StaffID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isMainAddress = false
				WHERE staffId = $1
			`, payload.StaffID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if payload.CustomerID != nil {
			_, err = tx.Exec(`
				UPDATE addresses
				SET isMainAddress = false
				WHERE customerId = $1
			`, payload.CustomerID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Update the address
	err = tx.QueryRow(`
		UPDATE addresses
		SET staffId = $1, customerId = $2, cityId = $3, street = $4, name = $5, postalCode = $6, buildingNumber = $7, isDeliveryAddress = $8, isMainAddress = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, staffId, customerId, cityId, street, name, postalCode, buildingNumber, isDeliveryAddress, isMainAddress, created_at, updated_at, deleted_at
	`, payload.StaffID, payload.CustomerID, payload.City.ID, payload.Street, payload.Name, payload.PostalCode, payload.BuildingNumber, payload.IsDeliveryAddress, payload.IsMainAddress, currentTime, foundAddress.ID).Scan(
		&address.ID,
		&address.StaffID,
		&address.CustomerID,
		&address.CityID,
		&address.Street,
		&address.Name,
		&address.PostalCode,
		&address.BuildingNumber,
		&address.IsDeliveryAddress,
		&address.IsMainAddress,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return the address
	return &address, nil
}
