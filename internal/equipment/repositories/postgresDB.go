package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/equipment/domain"
	"github.com/hoitek/Maja-Service/internal/equipment/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type EquipmentRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewEquipmentRepositoryPostgresDB(d *sql.DB) *EquipmentRepositoryPostgresDB {
	return &EquipmentRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.EquipmentsQueryRequestParams) []string {
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

func (r *EquipmentRepositoryPostgresDB) Query(queries *models.EquipmentsQueryRequestParams) ([]*domain.Equipment, error) {
	q := `SELECT * FROM equipments ls `
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

	var equipments []*domain.Equipment
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			equipment   domain.Equipment
			description sql.NullString
			deletedAt   pq.NullTime
		)
		err := rows.Scan(&equipment.ID, &equipment.Name, &description, &equipment.CreatedAt, &equipment.UpdatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			equipment.Description = &description.String
		}
		if deletedAt.Valid {
			equipment.DeletedAt = &deletedAt.Time
		}
		equipments = append(equipments, &equipment)
	}
	return equipments, nil
}

func (r *EquipmentRepositoryPostgresDB) Count(queries *models.EquipmentsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM equipments ls `
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

func (r *EquipmentRepositoryPostgresDB) Create(payload *models.EquipmentsCreateRequestBody) (*domain.Equipment, error) {
	var equipment domain.Equipment

	// Current time
	currentTime := time.Now()

	// Insert the equipment
	err := r.PostgresDB.QueryRow(`
		INSERT INTO equipments (name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&equipment.ID, &equipment.Name, &equipment.Description, &equipment.CreatedAt, &equipment.UpdatedAt, &equipment.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the equipment
	return &equipment, nil
}

func (r *EquipmentRepositoryPostgresDB) Delete(payload *models.EquipmentsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM equipments
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

func (r *EquipmentRepositoryPostgresDB) Update(payload *models.EquipmentsCreateRequestBody, id int64) (*domain.Equipment, error) {
	var equipment domain.Equipment

	// Current time
	currentTime := time.Now()

	// Find the equipment by name
	var foundEquipment domain.Equipment
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM equipments
		WHERE id = $1
	`, id).Scan(&foundEquipment.ID, &foundEquipment.Name, &foundEquipment.Description, &foundEquipment.CreatedAt, &foundEquipment.UpdatedAt, &foundEquipment.DeletedAt)

	// If the equipment is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the equipment
	err = r.PostgresDB.QueryRow(`
		UPDATE equipments
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, foundEquipment.ID).Scan(&equipment.ID, &equipment.Name, &equipment.Description, &equipment.CreatedAt, &equipment.UpdatedAt, &equipment.DeletedAt)

	// If the equipment does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the equipment
	return &equipment, nil
}

func (r *EquipmentRepositoryPostgresDB) GetEquipmentsByIds(ids []int64) ([]*domain.Equipment, error) {
	var equipments []*domain.Equipment
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM equipments
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			equipment   domain.Equipment
			description sql.NullString
			deletedAt   pq.NullTime
		)
		err := rows.Scan(&equipment.ID, &equipment.Name, &description, &equipment.CreatedAt, &equipment.UpdatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			equipment.Description = &description.String
		}
		if deletedAt.Valid {
			equipment.DeletedAt = &deletedAt.Time
		}
		equipments = append(equipments, &equipment)
	}
	return equipments, nil
}
