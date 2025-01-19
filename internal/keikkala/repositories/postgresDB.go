package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/keikkala/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/keikkala/domain"
	"github.com/hoitek/Maja-Service/internal/keikkala/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type KeikkalaRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewKeikkalaRepositoryPostgresDB(d *sql.DB) *KeikkalaRepositoryPostgresDB {
	return &KeikkalaRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.KeikkalasQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" k.id = %d", queries.ID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.SectionIDs.Op != "" {
			var sectionIDs []int
			err := json.Unmarshal([]byte(queries.Filters.SectionIDs.Value), &sectionIDs)
			if err != nil {
				log.Println("Error unmarshaling sectionIDs in keikkala repository filters", err.Error())
			} else {
				where = append(where, "EXISTS (SELECT 1 FROM jsonb_array_elements_text(k.sections) AS section WHERE section::int IN ("+strings.Join(strings.Fields(strings.Trim(fmt.Sprint(sectionIDs), "[]")), ",")+") )")
			}
		}
		if queries.Filters.RoleID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleID.Op, fmt.Sprintf("%v", queries.Filters.RoleID.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.roleId %s %s", opValue.Operator, val))
		}
		if queries.Filters.RoleName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.RoleName.Op, fmt.Sprintf("%v", queries.Filters.RoleName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.StartDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartDate.Op, fmt.Sprintf("%v", queries.Filters.StartDate.Value))
			if opValue.Value != "all" {
				val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
				where = append(where, fmt.Sprintf(" k.start_date %s %s", opValue.Operator, val))
			}
		}
		if queries.Filters.EndDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndDate.Op, fmt.Sprintf("%v", queries.Filters.EndDate.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.end_date %s %s", opValue.Operator, val))
		}
		if queries.Filters.StartTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartTime.Op, fmt.Sprintf("%v", queries.Filters.StartTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.start_time %s %s", opValue.Operator, val))
		}
		if queries.Filters.EndTime.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EndTime.Op, fmt.Sprintf("%v", queries.Filters.EndTime.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.end_time %s %s", opValue.Operator, val))
		}
		if queries.Filters.KaupunkiAddress.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.KaupunkiAddress.Op, fmt.Sprintf("%v", queries.Filters.KaupunkiAddress.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.kaupunkiAddress %s %s", opValue.Operator, val))
		}
		if queries.Filters.PaymentType.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PaymentType.Op, fmt.Sprintf("%v", queries.Filters.PaymentType.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.paymentType %s %s", opValue.Operator, val))
		}
		if queries.Filters.ShiftName.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.ShiftName.Op, fmt.Sprintf("%v", queries.Filters.ShiftName.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(` k.shiftName %s %s`, opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.description %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.status %s %s", opValue.Operator, val))
		}
		if queries.Filters.PickedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PickedAt.Op, fmt.Sprintf("%v", queries.Filters.PickedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.picked_at %s %s", opValue.Operator, val))
		}
		if queries.Filters.PickedBy.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.PickedBy.Op, fmt.Sprintf("%v", queries.Filters.PickedBy.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.pickedBy %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.created_at %s %s", opValue.Operator, val))
		}
		if queries.Filters.UpdatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.UpdatedAt.Op, fmt.Sprintf("%v", queries.Filters.UpdatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" k.updated_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *KeikkalaRepositoryPostgresDB) Query(queries *models.KeikkalasQueryRequestParams) ([]*domain.Keikkala, error) {
	q := `
		SELECT
		    k.id,
		    k.roleId,
			k.start_date,
			k.end_date,
			k.start_time,
			k.end_time,
			k.kaupunkiAddress,
			k.sections,
			k.paymentType,
			k.shiftName,
			k.description,
			k.status,
			k.picked_at,
			k.created_at,
			k.updated_at,
			k.deleted_at,
			r.id AS roleId,
			r.name AS roleName,
			u.id AS pickedById,
			u.firstName AS pickedByFirstName,
			u.lastName AS pickedByLastName,
			u.email AS pickedByEmail,
			u.avatarUrl AS pickedByAvatarUrl,
			u2.id AS createdById,
			u2.firstName AS createdByFirstName,
			u2.lastName AS createdByLastName,
			u2.email AS createdByEmail,
			u2.avatarUrl AS createdByAvatarUrl,
			u3.id AS updatedById,
			u3.firstName AS updatedByFirstName,
			u3.lastName AS updatedByLastName,
			u3.email AS updatedByEmail,
			u3.avatarUrl AS updatedByAvatarUrl,
			s2.sections AS joinedSections
		FROM keikkalaShifts k
		LEFT JOIN _roles r ON r.id = k.roleId
		LEFT JOIN users u ON u.id = k.pickedBy
		LEFT JOIN users u2 ON u2.id = k.createdBy
		LEFT JOIN users u3 ON u3.id = k.updatedBy
		LEFT JOIN LATERAL (
			SELECT jsonb_agg(s) AS sections
			FROM (
				SELECT
					s.id,
					s.name
				FROM sections s
				WHERE s.id::text IN (SELECT jsonb_array_elements_text(k.sections)::text)
			) s
		) AS s2 ON true
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" k.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.StartDate.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" k.start_date %s", queries.Sorts.StartDate.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" k.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var keikkalas []*domain.Keikkala
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			keikkala            domain.Keikkala
			roleId              sql.NullInt64
			kaupunkiAddress     sql.NullString
			sections            json.RawMessage
			description         sql.NullString
			pickedAt            pq.NullTime
			deletedAt           pq.NullTime
			kRoleId             sql.NullInt64
			kRoleName           sql.NullString
			kPickedById         sql.NullInt64
			kPickedByFirstName  sql.NullString
			kPickedByLastName   sql.NullString
			kPickedByEmail      sql.NullString
			kPickedByAvatarUrl  sql.NullString
			kCreatedById        sql.NullInt64
			kCreatedByFirstName sql.NullString
			kCreatedByLastName  sql.NullString
			kCreatedByEmail     sql.NullString
			kCreatedByAvatarUrl sql.NullString
			kUpdatedById        sql.NullInt64
			kUpdatedByFirstName sql.NullString
			kUpdatedByLastName  sql.NullString
			kUpdatedByEmail     sql.NullString
			kUpdatedByAvatarUrl sql.NullString
			joinedSections      json.RawMessage
		)
		err := rows.Scan(
			&keikkala.ID,
			&roleId,
			&keikkala.StartDate,
			&keikkala.EndDate,
			&keikkala.StartTime,
			&keikkala.EndTime,
			&kaupunkiAddress,
			&sections,
			&keikkala.PaymentType,
			&keikkala.ShiftName,
			&description,
			&keikkala.Status,
			&pickedAt,
			&keikkala.CreatedAt,
			&keikkala.UpdatedAt,
			&deletedAt,
			&kRoleId,
			&kRoleName,
			&kPickedById,
			&kPickedByFirstName,
			&kPickedByLastName,
			&kPickedByEmail,
			&kPickedByAvatarUrl,
			&kCreatedById,
			&kCreatedByFirstName,
			&kCreatedByLastName,
			&kCreatedByEmail,
			&kCreatedByAvatarUrl,
			&kUpdatedById,
			&kUpdatedByFirstName,
			&kUpdatedByLastName,
			&kUpdatedByEmail,
			&kUpdatedByAvatarUrl,
			&joinedSections,
		)
		if err != nil {
			return nil, err
		}
		if roleId.Valid {
			rid := uint(roleId.Int64)
			keikkala.RoleID = &rid
		}
		if kaupunkiAddress.Valid {
			keikkala.KaupunkiAddress = &kaupunkiAddress.String
		}
		if description.Valid {
			keikkala.Description = &description.String
		}
		if pickedAt.Valid {
			keikkala.PickedAt = &pickedAt.Time
		}
		if kPickedById.Valid {
			keikkala.PickedBy = &domain.KeikkalaUser{
				ID: uint(kPickedById.Int64),
			}
			if kPickedByFirstName.Valid {
				keikkala.PickedBy.FirstName = kPickedByFirstName.String
			}
			if kPickedByLastName.Valid {
				keikkala.PickedBy.LastName = kPickedByLastName.String
			}
			if kPickedByEmail.Valid {
				keikkala.PickedBy.Email = kPickedByEmail.String
			}
		}
		if deletedAt.Valid {
			keikkala.DeletedAt = &deletedAt.Time
		}
		if kCreatedById.Valid {
			keikkala.CreatedBy = &domain.KeikkalaUser{
				ID: uint(kCreatedById.Int64),
			}
			if kCreatedByFirstName.Valid {
				keikkala.CreatedBy.FirstName = kCreatedByFirstName.String
			}
			if kCreatedByLastName.Valid {
				keikkala.CreatedBy.LastName = kCreatedByLastName.String
			}
			if kCreatedByEmail.Valid {
				keikkala.CreatedBy.Email = kCreatedByEmail.String
			}
		}
		if kUpdatedById.Valid {
			keikkala.UpdatedBy = &domain.KeikkalaUser{
				ID: uint(kUpdatedById.Int64),
			}
			if kUpdatedByFirstName.Valid {
				keikkala.UpdatedBy.FirstName = kUpdatedByFirstName.String
			}
			if kUpdatedByLastName.Valid {
				keikkala.UpdatedBy.LastName = kUpdatedByLastName.String
			}
			if kUpdatedByEmail.Valid {
				keikkala.UpdatedBy.Email = kUpdatedByEmail.String
			}
		}
		if kRoleId.Valid {
			keikkala.Role = &domain.KeikkalaRole{
				ID:   uint(kRoleId.Int64),
				Name: kRoleName.String,
			}
		}
		if joinedSections != nil {
			var joinedSectionsData []*domain.KeikkalaSection
			err := json.Unmarshal(joinedSections, &joinedSectionsData)
			if err != nil {
				log.Printf("error unmarshalling joined sections: %v", err.Error())
			} else {
				keikkala.Sections = joinedSectionsData
			}
		}
		keikkalas = append(keikkalas, &keikkala)
	}
	return keikkalas, nil
}

func (r *KeikkalaRepositoryPostgresDB) Count(queries *models.KeikkalasQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(k.id)
		FROM keikkalaShifts k
		LEFT JOIN _roles r ON r.id = k.roleId
		LEFT JOIN users u ON u.id = k.pickedBy
		LEFT JOIN users u2 ON u2.id = k.createdBy
		LEFT JOIN users u3 ON u3.id = k.updatedBy
		LEFT JOIN LATERAL (
			SELECT jsonb_agg(s) AS sections
			FROM (
				SELECT
					s.id,
					s.name
				FROM sections s
				WHERE s.id::text IN (SELECT jsonb_array_elements_text(k.sections)::text)
			) s
		) AS s2 ON true
	`
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

func (r *KeikkalaRepositoryPostgresDB) Create(payload *models.KeikkalasCreateRequestBody) (*domain.Keikkala, error) {
	// Current time
	currentTime := time.Now()

	// Create sections jsonb
	sectionsBytes, err := json.Marshal(payload.SectionIDsInt64)
	if err != nil {
		return nil, err
	}
	sectionsJSON := string(sectionsBytes)

	// Create shift name
	var shiftName = "morning"
	if payload.StartTimeAsTime.Hour() >= 12 {
		shiftName = "evening"
	}
	if payload.StartTimeAsTime.Hour() >= 18 {
		shiftName = "night"
	}

	// Insert the keikkala
	var insertedId int
	err = r.PostgresDB.QueryRow(`
		INSERT INTO keikkalaShifts (
			start_date,
			end_date,
			start_time,
			end_time,
			roleId,
			kaupunkiAddress,
			sections,
			paymentType,
		    shiftName,
			description,
			createdBy,
			updatedBy,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.StartTimeAsTime,
		payload.EndTimeAsTime,
		payload.RoleID,
		payload.KaupunkiAddress,
		sectionsJSON,
		payload.PaymentType,
		shiftName,
		payload.Description,
		payload.AuthenticatedUser.ID,
		payload.AuthenticatedUser.ID,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the keikkala
	keikkalas, err := r.Query(&models.KeikkalasQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(keikkalas) == 0 {
		return nil, errors.New("no rows affected")
	}
	return keikkalas[0], nil
}

func (r *KeikkalaRepositoryPostgresDB) Delete(payload *models.KeikkalasDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM keikkalas
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

func (r *KeikkalaRepositoryPostgresDB) GetKeikkalaShiftsByIds(ids []int64) ([]*domain.Keikkala, error) {
	var keikkalas []*domain.Keikkala
	rows, err := r.PostgresDB.Query(`
		SELECT
		    k.id,
		    k.roleId,
			k.start_date,
			k.end_date,
			k.start_time,
			k.end_time,
			k.kaupunkiAddress,
			k.sections,
			k.paymentType,
			k.shiftName,
			k.description,
			k.status,
			k.picked_at,
			k.created_at,
			k.updated_at,
			k.deleted_at,
			r.id AS roleId,
			r.name AS roleName,
			u.id AS pickedById,
			u.firstName AS pickedByFirstName,
			u.lastName AS pickedByLastName,
			u.email AS pickedByEmail,
			u.avatarUrl AS pickedByAvatarUrl,
			u2.id AS createdById,
			u2.firstName AS createdByFirstName,
			u2.lastName AS createdByLastName,
			u2.email AS createdByEmail,
			u2.avatarUrl AS createdByAvatarUrl,
			u3.id AS updatedById,
			u3.firstName AS updatedByFirstName,
			u3.lastName AS updatedByLastName,
			u3.email AS updatedByEmail,
			u3.avatarUrl AS updatedByAvatarUrl,
			s2.sections AS joinedSections
		FROM keikkalaShifts k
		LEFT JOIN _roles r ON r.id = k.roleId
		LEFT JOIN users u ON u.id = k.pickedBy
		LEFT JOIN users u2 ON u2.id = k.createdBy
		LEFT JOIN users u3 ON u3.id = k.updatedBy
		LEFT JOIN LATERAL (
			SELECT jsonb_agg(s) AS sections
			FROM (
				SELECT
					s.id,
					s.name
				FROM sections s
				WHERE s.id::text IN (SELECT jsonb_array_elements_text(k.sections)::text)
			) s
		) AS s2 ON true
		WHERE k.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			keikkala            domain.Keikkala
			roleId              sql.NullInt64
			kaupunkiAddress     sql.NullString
			sections            json.RawMessage
			description         sql.NullString
			pickedAt            pq.NullTime
			deletedAt           pq.NullTime
			kRoleId             sql.NullInt64
			kRoleName           sql.NullString
			kPickedById         sql.NullInt64
			kPickedByFirstName  sql.NullString
			kPickedByLastName   sql.NullString
			kPickedByEmail      sql.NullString
			kPickedByAvatarUrl  sql.NullString
			kCreatedById        sql.NullInt64
			kCreatedByFirstName sql.NullString
			kCreatedByLastName  sql.NullString
			kCreatedByEmail     sql.NullString
			kCreatedByAvatarUrl sql.NullString
			kUpdatedById        sql.NullInt64
			kUpdatedByFirstName sql.NullString
			kUpdatedByLastName  sql.NullString
			kUpdatedByEmail     sql.NullString
			kUpdatedByAvatarUrl sql.NullString
			joinedSections      json.RawMessage
		)
		err := rows.Scan(
			&keikkala.ID,
			&roleId,
			&keikkala.StartDate,
			&keikkala.EndDate,
			&keikkala.StartTime,
			&keikkala.EndTime,
			&kaupunkiAddress,
			&sections,
			&keikkala.PaymentType,
			&keikkala.ShiftName,
			&description,
			&keikkala.Status,
			&pickedAt,
			&keikkala.CreatedAt,
			&keikkala.UpdatedAt,
			&deletedAt,
			&kRoleId,
			&kRoleName,
			&kPickedById,
			&kPickedByFirstName,
			&kPickedByLastName,
			&kPickedByEmail,
			&kPickedByAvatarUrl,
			&kCreatedById,
			&kCreatedByFirstName,
			&kCreatedByLastName,
			&kCreatedByEmail,
			&kCreatedByAvatarUrl,
			&kUpdatedById,
			&kUpdatedByFirstName,
			&kUpdatedByLastName,
			&kUpdatedByEmail,
			&kUpdatedByAvatarUrl,
			&joinedSections,
		)
		if err != nil {
			return nil, err
		}
		if roleId.Valid {
			rid := uint(roleId.Int64)
			keikkala.RoleID = &rid
		}
		if kaupunkiAddress.Valid {
			keikkala.KaupunkiAddress = &kaupunkiAddress.String
		}
		if description.Valid {
			keikkala.Description = &description.String
		}
		if pickedAt.Valid {
			keikkala.PickedAt = &pickedAt.Time
		}
		if kPickedById.Valid {
			keikkala.PickedBy = &domain.KeikkalaUser{
				ID: uint(kPickedById.Int64),
			}
			if kPickedByFirstName.Valid {
				keikkala.PickedBy.FirstName = kPickedByFirstName.String
			}
			if kPickedByLastName.Valid {
				keikkala.PickedBy.LastName = kPickedByLastName.String
			}
			if kPickedByEmail.Valid {
				keikkala.PickedBy.Email = kPickedByEmail.String
			}
		}
		if deletedAt.Valid {
			keikkala.DeletedAt = &deletedAt.Time
		}
		if kCreatedById.Valid {
			keikkala.CreatedBy = &domain.KeikkalaUser{
				ID: uint(kCreatedById.Int64),
			}
			if kCreatedByFirstName.Valid {
				keikkala.CreatedBy.FirstName = kCreatedByFirstName.String
			}
			if kCreatedByLastName.Valid {
				keikkala.CreatedBy.LastName = kCreatedByLastName.String
			}
			if kCreatedByEmail.Valid {
				keikkala.CreatedBy.Email = kCreatedByEmail.String
			}
		}
		if kUpdatedById.Valid {
			keikkala.UpdatedBy = &domain.KeikkalaUser{
				ID: uint(kUpdatedById.Int64),
			}
			if kUpdatedByFirstName.Valid {
				keikkala.UpdatedBy.FirstName = kUpdatedByFirstName.String
			}
			if kUpdatedByLastName.Valid {
				keikkala.UpdatedBy.LastName = kUpdatedByLastName.String
			}
			if kUpdatedByEmail.Valid {
				keikkala.UpdatedBy.Email = kUpdatedByEmail.String
			}
		}
		if kRoleId.Valid {
			keikkala.Role = &domain.KeikkalaRole{
				ID:   uint(kRoleId.Int64),
				Name: kRoleName.String,
			}
		}
		if joinedSections != nil {
			var joinedSectionsData []*domain.KeikkalaSection
			err := json.Unmarshal(joinedSections, &joinedSectionsData)
			if err != nil {
				log.Printf("error unmarshalling joined sections: %v", err.Error())
			} else {
				keikkala.Sections = joinedSectionsData
			}
		}
		keikkalas = append(keikkalas, &keikkala)
	}
	return keikkalas, nil
}

func (r *KeikkalaRepositoryPostgresDB) QueryShiftStatistics(queries *models.KeikkalasQueryShiftStatisticsRequestParams) (int, int, int, error) {
	var q = `
		SELECT
			COUNT(k.id)
		FROM keikkalaShifts k
	`
	var (
		morningCount int64
		eveningCount int64
		nightCount   int64
	)
	for _, shiftName := range []string{constants.KEIKKALA_SHIFT_NAME_MORNING, constants.KEIKKALA_SHIFT_NAME_EVENING, constants.KEIKKALA_SHIFT_NAME_NIGHT} {
		var count int64
		qShift := fmt.Sprintf("%s WHERE k.shiftName = '%s' AND k.deleted_at IS NULL", q, shiftName)
		if queries.Filters.StartDate.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StartDate.Op, fmt.Sprintf("%v", queries.Filters.StartDate.Value))
			if opValue.Value != "all" {
				val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
				qShift = fmt.Sprintf("%s AND k.start_date %s %s", qShift, opValue.Operator, val)
			}
		}
		err := r.PostgresDB.QueryRow(qShift).Scan(&count)
		if err != nil {
			return 0, 0, 0, err
		}
		switch shiftName {
		case constants.KEIKKALA_SHIFT_NAME_MORNING:
			morningCount = count
		case constants.KEIKKALA_SHIFT_NAME_EVENING:
			eveningCount = count
		case constants.KEIKKALA_SHIFT_NAME_NIGHT:
			nightCount = count
		}
	}
	return int(morningCount), int(eveningCount), int(nightCount), nil
}
