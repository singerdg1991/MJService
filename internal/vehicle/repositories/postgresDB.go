package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
	"github.com/hoitek/Maja-Service/internal/vehicle/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type VehicleRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewVehicleRepositoryPostgresDB(d *sql.DB) *VehicleRepositoryPostgresDB {
	return &VehicleRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.VehiclesQueryRequestParams) []string {
	var where []string
	if queries.ID > 0 {
		where = append(where, fmt.Sprintf("vehicles.id = %d", queries.ID))
	}
	if queries.Filters.UserID.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.UserID.Op, fmt.Sprintf("%v", queries.Filters.UserID.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.userId %s %s", opValue.Operator, val))
	}
	if queries.Filters.VehicleType.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.VehicleType.Op, fmt.Sprintf("%v", queries.Filters.VehicleType.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.vehicleType %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.Brand.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.Brand.Op, fmt.Sprintf("%v", queries.Filters.Brand.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.brand %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.Model.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.Model.Op, fmt.Sprintf("%v", queries.Filters.Model.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.model %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.Year.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.Year.Op, fmt.Sprintf("%v", queries.Filters.Year.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.year %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.Variant.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.Variant.Op, fmt.Sprintf("%v", queries.Filters.Variant.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.variant %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.FuelType.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.FuelType.Op, fmt.Sprintf("%v", queries.Filters.FuelType.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.fuelType %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.VehicleNo.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.VehicleNo.Op, fmt.Sprintf("%v", queries.Filters.VehicleNo.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.vehicleNo %s '%s'", opValue.Operator, val))
	}
	if queries.Filters.CreatedAt.Op != "" {
		opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("vehicles.created_at %s '%s'", opValue.Operator, val))
	}
	return where
}

func (r *VehicleRepositoryPostgresDB) Query(queries *models.VehiclesQueryRequestParams) ([]*domain.Vehicle, error) {
	q := `
		SELECT
			vehicles.*,
			users.id AS ownerUserId,
			users.firstName AS ownerUserFirstName,
			users.lastName AS ownerUserLastName,
			users.avatarUrl AS ownerUserAvatarUrl
		FROM vehicles
		LEFT JOIN users ON users.id = vehicles.userId
	`

	// Add where clause
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += fmt.Sprintf("WHERE %s ", strings.Join(where, "AND"))
		}
		limit := exp.TerIf(queries.Limit == 0, 10, queries.Limit)
		queries.Page = exp.TerIf(queries.Page == 0, 1, queries.Page)
		offset := (queries.Page - 1) * limit
		q += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	// Execute query
	var vehicles []*domain.Vehicle
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			vehicle            domain.Vehicle
			ownerUserId        sql.NullInt64
			ownerUserFirstName sql.NullString
			ownerUserLastName  sql.NullString
			ownerUserAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&vehicle.ID,
			&vehicle.VehicleType,
			&vehicle.UserID,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.Year,
			&vehicle.Variant,
			&vehicle.FuelType,
			&vehicle.VehicleNo,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
			&vehicle.DeletedAt,
			&ownerUserId,
			&ownerUserFirstName,
			&ownerUserLastName,
			&ownerUserAvatarUrl,
		)
		if err != nil {
			return nil, err
		}

		// Owner
		vehicle.User = &domain.VehicleUser{}
		if ownerUserId.Valid {
			vehicle.User.ID = uint(ownerUserId.Int64)
		}
		if ownerUserFirstName.Valid {
			vehicle.User.FirstName = ownerUserFirstName.String
		}
		if ownerUserLastName.Valid {
			vehicle.User.LastName = ownerUserLastName.String
		}
		if ownerUserAvatarUrl.Valid {
			vehicle.User.AvatarUrl = ownerUserAvatarUrl.String
		}
		if vehicle.User.ID == 0 {
			vehicle.User = nil
		}

		// Append to slice
		vehicles = append(vehicles, &vehicle)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepositoryPostgresDB) Count(queries *models.VehiclesQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    count(vehicles.id)
		FROM vehicles
		LEFT JOIN users ON users.id = vehicles.userId
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += fmt.Sprintf("WHERE %s ", strings.Join(where, "AND"))
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *VehicleRepositoryPostgresDB) Create(payload *models.VehiclesCreateRequestBody) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle

	// Current time
	currentTime := time.Now()

	// Insert the vehicle
	err := r.PostgresDB.QueryRow(`
	INSERT INTO vehicles (vehicleType, userId, brand, model, year, variant, fuelType, vehicleNo, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, vehicleType, userId, brand, model, year, variant, fuelType, vehicleNo, created_at, updated_at, deleted_at
`, payload.VehicleType, payload.UserID, payload.Brand, payload.Model, payload.Year, payload.Variant, payload.FuelType, payload.VehicleNo, currentTime, currentTime, nil).Scan(
		&vehicle.ID,
		&vehicle.VehicleType,
		&vehicle.UserID,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.Variant,
		&vehicle.FuelType,
		&vehicle.VehicleNo,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
		&vehicle.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *VehicleRepositoryPostgresDB) Delete(payload *models.VehiclesDeleteRequestBody) ([]int64, error) {
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
		DELETE FROM vehicles
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

func (r *VehicleRepositoryPostgresDB) Update(payload *models.VehiclesCreateRequestBody, id int) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle

	// Current time
	currentTime := time.Now()

	// Update the vehicle
	err := r.PostgresDB.QueryRow(`
		UPDATE vehicles
		SET vehicleType = $1, userId = $2, brand = $3, year = $4, variant = $5, fuelType = $6, vehicleNo = $7, updated_at = $8
		WHERE id = $9
		RETURNING id, vehicleType, userId, brand, year, variant, fuelType, vehicleNo, created_at, updated_at, deleted_at
	`, payload.VehicleType, payload.UserID, payload.Brand, payload.Year, payload.Variant, payload.FuelType, payload.VehicleNo, currentTime, id).Scan(
		&vehicle.ID,
		&vehicle.VehicleType,
		&vehicle.UserID,
		&vehicle.Brand,
		&vehicle.Year,
		&vehicle.Variant,
		&vehicle.FuelType,
		&vehicle.VehicleNo,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
		&vehicle.DeletedAt,
	)

	// If the vehicle does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the vehicle
	return &vehicle, nil
}
