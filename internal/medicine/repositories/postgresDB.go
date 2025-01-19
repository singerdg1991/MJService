package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/medicine/domain"
	"github.com/hoitek/Maja-Service/internal/medicine/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type MedicineRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewMedicineRepositoryPostgresDB(d *sql.DB) *MedicineRepositoryPostgresDB {
	return &MedicineRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.MedicinesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" m.id = %d", queries.ID))
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.Code.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Code.Op, fmt.Sprintf("%v", queries.Filters.Code.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.code %s %s", opValue.Operator, val))
		}
		if queries.Filters.Availability.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Availability.Op, fmt.Sprintf("%v", queries.Filters.Availability.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.availability %s %s", opValue.Operator, val))
		}
		if queries.Filters.Manufacturer.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Manufacturer.Op, fmt.Sprintf("%v", queries.Filters.Manufacturer.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.manufacturer %s %s", opValue.Operator, val))
		}
		if queries.Filters.PurposeOfUse.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PurposeOfUse.Op, fmt.Sprintf("%v", queries.Filters.PurposeOfUse.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.purposeOfUse %s %s", opValue.Operator, val))
		}
		if queries.Filters.Instruction.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Instruction.Op, fmt.Sprintf("%v", queries.Filters.Instruction.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.instruction %s %s", opValue.Operator, val))
		}
		if queries.Filters.SideEffects.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.SideEffects.Op, fmt.Sprintf("%v", queries.Filters.SideEffects.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.sideEffects %s %s", opValue.Operator, val))
		}
		if queries.Filters.Conditions.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Conditions.Op, fmt.Sprintf("%v", queries.Filters.Conditions.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.conditions %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" m.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *MedicineRepositoryPostgresDB) Query(queries *models.MedicinesQueryRequestParams) ([]*domain.Medicine, error) {
	q := `SELECT * FROM medicines m `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Code.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.code %s", queries.Sorts.Code.Op))
		}
		if queries.Sorts.Availability.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.availability %s", queries.Sorts.Availability.Op))
		}
		if queries.Sorts.Manufacturer.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.manufacturer %s", queries.Sorts.Manufacturer.Op))
		}
		if queries.Sorts.PurposeOfUse.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.purposeOfUse %s", queries.Sorts.PurposeOfUse.Op))
		}
		if queries.Sorts.Instruction.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.instruction %s", queries.Sorts.Instruction.Op))
		}
		if queries.Sorts.SideEffects.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.sideEffects %s", queries.Sorts.SideEffects.Op))
		}
		if queries.Sorts.Conditions.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.conditions %s", queries.Sorts.Conditions.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.description %s", queries.Sorts.Description.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" m.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var medicines []*domain.Medicine
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			medicine     domain.Medicine
			code         sql.NullString
			availability sql.NullString
			manufacturer sql.NullString
			purposeOfUse sql.NullString
			instruction  sql.NullString
			sideEffects  sql.NullString
			conditions   sql.NullString
			description  sql.NullString
			deletedAt    sql.NullTime
		)
		err := rows.Scan(
			&medicine.ID,
			&medicine.Name,
			&code,
			&availability,
			&manufacturer,
			&purposeOfUse,
			&instruction,
			&sideEffects,
			&conditions,
			&description,
			&medicine.CreatedAt,
			&medicine.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if code.Valid {
			medicine.Code = &code.String
		}
		if availability.Valid {
			medicine.Availability = &availability.String
		}
		if manufacturer.Valid {
			medicine.Manufacturer = &manufacturer.String
		}
		if purposeOfUse.Valid {
			medicine.PurposeOfUse = &purposeOfUse.String
		}
		if instruction.Valid {
			medicine.Instruction = &instruction.String
		}
		if sideEffects.Valid {
			medicine.SideEffects = &sideEffects.String
		}
		if conditions.Valid {
			medicine.Conditions = &conditions.String
		}
		if description.Valid {
			medicine.Description = &description.String
		}
		if deletedAt.Valid {
			medicine.DeletedAt = &deletedAt.Time
		}
		medicines = append(medicines, &medicine)
	}
	return medicines, nil
}

func (r *MedicineRepositoryPostgresDB) Count(queries *models.MedicinesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(m.id) FROM medicines m `
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

func (r *MedicineRepositoryPostgresDB) Create(payload *models.MedicinesCreateRequestBody) (*domain.Medicine, error) {
	var (
		medicine     domain.Medicine
		code         sql.NullString
		availability sql.NullString
		manufacturer sql.NullString
		purposeOfUse sql.NullString
		instruction  sql.NullString
		sideEffects  sql.NullString
		conditions   sql.NullString
		description  sql.NullString
		deletedAt    sql.NullTime
	)

	// Current time
	currentTime := time.Now()

	// Insert the medicine
	err := r.PostgresDB.QueryRow(`
		INSERT INTO medicines (name, code, availability, manufacturer, purposeOfUse, instruction, sideEffects, conditions, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, code, availability, manufacturer, purposeOfUse, instruction, sideEffects, conditions, description, created_at, updated_at, deleted_at
	`,
		payload.Name,
		payload.Code,
		payload.Availability,
		payload.Manufacturer,
		payload.PurposeOfUse,
		payload.Instruction,
		payload.SideEffects,
		payload.Conditions,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(
		&medicine.ID,
		&medicine.Name,
		&code,
		&availability,
		&manufacturer,
		&purposeOfUse,
		&instruction,
		&sideEffects,
		&conditions,
		&description,
		&medicine.CreatedAt,
		&medicine.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}
	if code.Valid {
		medicine.Code = &code.String
	}
	if availability.Valid {
		medicine.Availability = &availability.String
	}
	if manufacturer.Valid {
		medicine.Manufacturer = &manufacturer.String
	}
	if purposeOfUse.Valid {
		medicine.PurposeOfUse = &purposeOfUse.String
	}
	if instruction.Valid {
		medicine.Instruction = &instruction.String
	}
	if sideEffects.Valid {
		medicine.SideEffects = &sideEffects.String
	}
	if conditions.Valid {
		medicine.Conditions = &conditions.String
	}
	if description.Valid {
		medicine.Description = &description.String
	}
	if deletedAt.Valid {
		medicine.DeletedAt = &deletedAt.Time
	}

	// Return the medicine
	return &medicine, nil
}

func (r *MedicineRepositoryPostgresDB) Delete(payload *models.MedicinesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM medicines
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

func (r *MedicineRepositoryPostgresDB) Update(payload *models.MedicinesCreateRequestBody, id int64) (*domain.Medicine, error) {
	var (
		medicine     domain.Medicine
		code         sql.NullString
		availability sql.NullString
		manufacturer sql.NullString
		purposeOfUse sql.NullString
		instruction  sql.NullString
		sideEffects  sql.NullString
		conditions   sql.NullString
		description  sql.NullString
		deletedAt    sql.NullTime
	)

	// Current time
	currentTime := time.Now()

	// Find the medicine by name
	var foundMedicine domain.Medicine
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM medicines
		WHERE id = $1
	`, id).Scan(
		&foundMedicine.ID,
		&foundMedicine.Name,
		&code,
		&availability,
		&manufacturer,
		&purposeOfUse,
		&instruction,
		&sideEffects,
		&conditions,
		&description,
		&medicine.CreatedAt,
		&medicine.UpdatedAt,
		&deletedAt,
	)

	// If the medicine is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	if code.Valid {
		foundMedicine.Code = &code.String
	}
	if availability.Valid {
		foundMedicine.Availability = &availability.String
	}
	if manufacturer.Valid {
		foundMedicine.Manufacturer = &manufacturer.String
	}
	if purposeOfUse.Valid {
		foundMedicine.PurposeOfUse = &purposeOfUse.String
	}
	if instruction.Valid {
		foundMedicine.Instruction = &instruction.String
	}
	if sideEffects.Valid {
		foundMedicine.SideEffects = &sideEffects.String
	}
	if conditions.Valid {
		foundMedicine.Conditions = &conditions.String
	}
	if description.Valid {
		foundMedicine.Description = &description.String
	}
	if deletedAt.Valid {
		foundMedicine.DeletedAt = &deletedAt.Time
	}

	// Update the medicine
	err = r.PostgresDB.QueryRow(`
		UPDATE medicines
		SET name = $1, code = $2, availability = $3, manufacturer = $4, purposeOfUse = $5, instruction = $6, sideEffects = $7, conditions = $8, description = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, name, code, availability, manufacturer, purposeOfUse, instruction, sideEffects, conditions, description, created_at, updated_at, deleted_at
	`,
		payload.Name,
		payload.Code,
		payload.Availability,
		payload.Manufacturer,
		payload.PurposeOfUse,
		payload.Instruction,
		payload.SideEffects,
		payload.Conditions,
		payload.Description,
		currentTime,
		foundMedicine.ID,
	).Scan(
		&medicine.ID,
		&medicine.Name,
		&code,
		&availability,
		&manufacturer,
		&purposeOfUse,
		&instruction,
		&sideEffects,
		&conditions,
		&description,
		&medicine.CreatedAt,
		&medicine.UpdatedAt,
		&deletedAt,
	)

	// If the medicine does not update, return an error
	if err != nil {
		return nil, err
	}

	if code.Valid {
		medicine.Code = &code.String
	}
	if availability.Valid {
		medicine.Availability = &availability.String
	}
	if manufacturer.Valid {
		medicine.Manufacturer = &manufacturer.String
	}
	if purposeOfUse.Valid {
		medicine.PurposeOfUse = &purposeOfUse.String
	}
	if instruction.Valid {
		medicine.Instruction = &instruction.String
	}
	if sideEffects.Valid {
		medicine.SideEffects = &sideEffects.String
	}
	if conditions.Valid {
		medicine.Conditions = &conditions.String
	}
	if description.Valid {
		medicine.Description = &description.String
	}
	if deletedAt.Valid {
		medicine.DeletedAt = &deletedAt.Time
	}

	// Return the medicine
	return &medicine, nil
}

func (r *MedicineRepositoryPostgresDB) GetMedicinesByIds(ids []int64) ([]*domain.Medicine, error) {
	var medicines []*domain.Medicine
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM medicines
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			medicine     domain.Medicine
			code         sql.NullString
			availability sql.NullString
			manufacturer sql.NullString
			purposeOfUse sql.NullString
			instruction  sql.NullString
			sideEffects  sql.NullString
			conditions   sql.NullString
			description  sql.NullString
			deletedAt    sql.NullTime
		)
		err := rows.Scan(
			&medicine.ID,
			&medicine.Name,
			&code,
			&availability,
			&manufacturer,
			&purposeOfUse,
			&instruction,
			&sideEffects,
			&conditions,
			&description,
			&medicine.CreatedAt,
			&medicine.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if code.Valid {
			medicine.Code = &code.String
		}
		if availability.Valid {
			medicine.Availability = &availability.String
		}
		if manufacturer.Valid {
			medicine.Manufacturer = &manufacturer.String
		}
		if purposeOfUse.Valid {
			medicine.PurposeOfUse = &purposeOfUse.String
		}
		if instruction.Valid {
			medicine.Instruction = &instruction.String
		}
		if sideEffects.Valid {
			medicine.SideEffects = &sideEffects.String
		}
		if conditions.Valid {
			medicine.Conditions = &conditions.String
		}
		if description.Valid {
			medicine.Description = &description.String
		}
		if deletedAt.Valid {
			medicine.DeletedAt = &deletedAt.Time
		}
		medicines = append(medicines, &medicine)
	}
	return medicines, nil
}
