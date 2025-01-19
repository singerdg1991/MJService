package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type HolidayRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewHolidayRepositoryPostgresDB(d *sql.DB) *HolidayRepositoryPostgresDB {
	return &HolidayRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.HolidaysQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" h.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" h.createdBy = %d", queries.UserID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" h.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" h.status %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *HolidayRepositoryPostgresDB) Query(queries *models.HolidaysQueryRequestParams) ([]*domain.Holiday, error) {
	q := `
		WITH RankedRoles AS (
			SELECT
				ur.userId,
				r.id AS roleId,
				r.name AS roleName,
				ROW_NUMBER() OVER (PARTITION BY ur.userId ORDER BY r.id) AS roleRank
			FROM usersRoles ur
			LEFT JOIN _roles r ON r.id = ur.roleId
		)
		SELECT
			h.id,
			h.start_date,
			h.end_date,
			h.title,
			h.paymentType,
			h.description,
			h.status,
			h.rejectedReason,
			h.accepted_at,
			h.rejected_at,
			h.created_at,
			h.updated_at,
			h.deleted_at,
			h.createdBy,
			h.updatedBy,
			u.id AS createdByUserId,
			u.firstName AS createdByFirstName,
			u.lastName AS createdByLastName,
			u.email AS createdByEmail,
			u.avatarUrl AS createdByAvatarUrl,
			u2.id AS updatedByUserId,
			u2.firstName AS updatedByFirstName,
			u2.lastName AS updatedByLastName,
			u2.email AS updatedByEmail,
			u2.avatarUrl AS updatedByAvatarUrl,
			r.id AS createdByRoleId,
			r.name AS createdByRoleName,
			r2.id AS updatedByRoleId,
			r2.name AS updatedByRoleName
		FROM staffClubHolidays h
		LEFT JOIN users u ON u.id = h.createdBy
		LEFT JOIN users u2 ON u2.id = h.updatedBy
		LEFT JOIN RankedRoles ur ON ur.userId = u.id AND ur.roleRank = 1
		LEFT JOIN _roles r ON r.id = ur.roleId
		LEFT JOIN RankedRoles ur2 ON ur2.userId = u2.id AND ur2.roleRank = 1
		LEFT JOIN _roles r2 ON r2.id = ur2.roleId
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" h.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" h.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" h.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var holidays []*domain.Holiday
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			holiday            domain.Holiday
			description        sql.NullString
			deletedAt          pq.NullTime
			createdByID        sql.NullInt64
			createdByFirstName sql.NullString
			createdByLastName  sql.NullString
			createdByEmail     sql.NullString
			createdByAvatarUrl sql.NullString
			updatedByID        sql.NullInt64
			updatedByFirstName sql.NullString
			updatedByLastName  sql.NullString
			updatedByEmail     sql.NullString
			updatedByAvatarUrl sql.NullString
			rejectedReason     sql.NullString
			acceptedAt         pq.NullTime
			rejectedAt         pq.NullTime
			createdBy          sql.NullInt64
			updatedBy          sql.NullInt64
			createdByRoleId    sql.NullInt64
			createdByRoleName  sql.NullString
			updatedByRoleId    sql.NullInt64
			updatedByRoleName  sql.NullString
		)
		err := rows.Scan(
			&holiday.ID,
			&holiday.StartDate,
			&holiday.EndDate,
			&holiday.Title,
			&holiday.PaymentType,
			&description,
			&holiday.Status,
			&rejectedReason,
			&acceptedAt,
			&rejectedAt,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
			&deletedAt,
			&createdBy,
			&updatedBy,
			&createdByID,
			&createdByFirstName,
			&createdByLastName,
			&createdByEmail,
			&createdByAvatarUrl,
			&updatedByID,
			&updatedByFirstName,
			&updatedByLastName,
			&updatedByEmail,
			&updatedByAvatarUrl,
			&createdByRoleId,
			&createdByRoleName,
			&updatedByRoleId,
			&updatedByRoleName,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			holiday.Description = &description.String
		}
		if deletedAt.Valid {
			holiday.DeletedAt = &deletedAt.Time
		}
		if rejectedReason.Valid {
			holiday.RejectedReason = &rejectedReason.String
		}
		if acceptedAt.Valid {
			holiday.AcceptedAt = &acceptedAt.Time
		}
		if rejectedAt.Valid {
			holiday.RejectedAt = &rejectedAt.Time
		}
		if createdByID.Valid {
			holiday.CreatedBy = domain.HolidayUser{
				ID: uint(createdByID.Int64),
			}
			if createdByFirstName.Valid {
				holiday.CreatedBy.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				holiday.CreatedBy.LastName = createdByLastName.String
			}
			if createdByEmail.Valid {
				holiday.CreatedBy.Email = createdByEmail.String
			}
			if createdByAvatarUrl.Valid {
				holiday.CreatedBy.AvatarUrl = createdByAvatarUrl.String
			}
		}
		if updatedByID.Valid {
			holiday.UpdatedBy = domain.HolidayUser{
				ID: uint(updatedByID.Int64),
			}
			if updatedByFirstName.Valid {
				holiday.UpdatedBy.FirstName = updatedByFirstName.String
			}
			if updatedByLastName.Valid {
				holiday.UpdatedBy.LastName = updatedByLastName.String
			}
			if updatedByEmail.Valid {
				holiday.UpdatedBy.Email = updatedByEmail.String
			}
			if updatedByAvatarUrl.Valid {
				holiday.UpdatedBy.AvatarUrl = updatedByAvatarUrl.String
			}
		}
		if createdByRoleId.Valid {
			holiday.CreatedBy.Role = &domain.HolidayUserRole{
				ID: uint(createdByRoleId.Int64),
			}
			if createdByRoleName.Valid {
				holiday.CreatedBy.Role.Name = createdByRoleName.String
			}
		}
		if updatedByRoleId.Valid {
			holiday.UpdatedBy.Role = &domain.HolidayUserRole{
				ID: uint(updatedByRoleId.Int64),
			}
			if updatedByRoleName.Valid {
				holiday.UpdatedBy.Role.Name = updatedByRoleName.String
			}
		}
		holidays = append(holidays, &holiday)
	}
	return holidays, nil
}

func (r *HolidayRepositoryPostgresDB) Count(queries *models.HolidaysQueryRequestParams) (int64, error) {
	q := `
		WITH RankedRoles AS (
			SELECT
				ur.userId,
				r.id AS roleId,
				r.name AS roleName,
				ROW_NUMBER() OVER (PARTITION BY ur.userId ORDER BY r.id) AS roleRank
			FROM usersRoles ur
			LEFT JOIN _roles r ON r.id = ur.roleId
		)
		SELECT
			COUNT(h.id)
		FROM staffClubHolidays h
		LEFT JOIN users u ON u.id = h.createdBy
		LEFT JOIN users u2 ON u2.id = h.updatedBy
		LEFT JOIN RankedRoles ur ON ur.userId = u.id AND ur.roleRank = 1
		LEFT JOIN _roles r ON r.id = ur.roleId
		LEFT JOIN RankedRoles ur2 ON ur2.userId = u2.id AND ur2.roleRank = 1
		LEFT JOIN _roles r2 ON r2.id = ur2.roleId
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

func (r *HolidayRepositoryPostgresDB) Create(payload *models.HolidaysCreateRequestBody) (*domain.Holiday, error) {
	// Current time
	currentTime := time.Now()

	// Insert the holiday
	var (
		insertedId  int
		status      = constants.HOLIDAY_STATUS_PENDING
		createdByID = payload.AuthenticatedUser.ID
	)
	if payload.User != nil {
		createdByID = payload.User.ID
	}
	err := r.PostgresDB.QueryRow(`
		INSERT INTO staffClubHolidays (start_date, end_date, title, paymentType, description, status, created_at, updated_at, createdBy, updatedBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Title,
		payload.PaymentType,
		payload.Description,
		status,
		currentTime,
		currentTime,
		createdByID,
		createdByID,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the holiday
	holidays, err := r.Query(&models.HolidaysQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(holidays) == 0 {
		return nil, errors.New("no rows affected")
	}
	return holidays[0], nil
}

func (r *HolidayRepositoryPostgresDB) Delete(payload *models.HolidaysDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM staffClubHolidays
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

func (r *HolidayRepositoryPostgresDB) Update(payload *models.HolidaysCreateRequestBody, id int64) (*domain.Holiday, error) {
	// Current time
	currentTime := time.Now()

	// Find the holiday by id
	foundHolidays, _ := r.Query(&models.HolidaysQueryRequestParams{ID: int(id)})
	if len(foundHolidays) == 0 {
		return nil, errors.New("holiday not found")
	}
	foundHoliday := foundHolidays[0]
	if foundHoliday.DeletedAt != nil {
		return nil, errors.New("holiday is deleted")
	}

	// Update the holiday
	var (
		updatedId   int
		createdByID = payload.AuthenticatedUser.ID
	)
	if payload.User != nil {
		createdByID = payload.User.ID
	}
	err := r.PostgresDB.QueryRow(`
		UPDATE staffClubHolidays
		SET start_date = $1, end_date = $2, title = $3, paymentType = $4, description = $5, updated_at = $6, updatedBy = $7, createdBy = $8
		WHERE id = $9
		RETURNING id
	`,
		payload.StartDateAsDate,
		payload.EndDateAsDate,
		payload.Title,
		payload.PaymentType,
		payload.Description,
		currentTime,
		createdByID,
		createdByID,
		id,
	).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the holiday
	holidays, err := r.Query(&models.HolidaysQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(holidays) == 0 {
		return nil, errors.New("no rows affected")
	}
	return holidays[0], nil
}

func (r *HolidayRepositoryPostgresDB) GetHolidaysByIds(ids []int64) ([]*domain.Holiday, error) {
	var holidays []*domain.Holiday
	rows, err := r.PostgresDB.Query(`
		WITH RankedRoles AS (
			SELECT
				ur.userId,
				r.id AS roleId,
				r.name AS roleName,
				ROW_NUMBER() OVER (PARTITION BY ur.userId ORDER BY r.id) AS roleRank
			FROM usersRoles ur
			LEFT JOIN _roles r ON r.id = ur.roleId
		)
		SELECT
			h.id,
			h.start_date,
			h.end_date,
			h.title,
			h.paymentType,
			h.description,
			h.status,
			h.rejectedReason,
			h.accepted_at,
			h.rejected_at,
			h.created_at,
			h.updated_at,
			h.deleted_at,
			h.createdBy,
			h.updatedBy,
			u.id AS createdByUserId,
			u.firstName AS createdByFirstName,
			u.lastName AS createdByLastName,
			u.email AS createdByEmail,
			u.avatarUrl AS createdByAvatarUrl,
			u2.id AS updatedByUserId,
			u2.firstName AS updatedByFirstName,
			u2.lastName AS updatedByLastName,
			u2.email AS updatedByEmail,
			u2.avatarUrl AS updatedByAvatarUrl,
			r.id AS createdByRoleId,
			r.name AS createdByRoleName,
			r2.id AS updatedByRoleId,
			r2.name AS updatedByRoleName
		FROM staffClubHolidays h
		LEFT JOIN users u ON u.id = h.createdBy
		LEFT JOIN users u2 ON u2.id = h.updatedBy
		LEFT JOIN RankedRoles ur ON ur.userId = u.id AND ur.roleRank = 1
		LEFT JOIN _roles r ON r.id = ur.roleId
		LEFT JOIN RankedRoles ur2 ON ur2.userId = u2.id AND ur2.roleRank = 1
		LEFT JOIN _roles r2 ON r2.id = ur2.roleId
		WHERE h.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			holiday            domain.Holiday
			description        sql.NullString
			deletedAt          pq.NullTime
			createdByID        sql.NullInt64
			createdByFirstName sql.NullString
			createdByLastName  sql.NullString
			createdByEmail     sql.NullString
			createdByAvatarUrl sql.NullString
			updatedByID        sql.NullInt64
			updatedByFirstName sql.NullString
			updatedByLastName  sql.NullString
			updatedByEmail     sql.NullString
			updatedByAvatarUrl sql.NullString
			rejectedReason     sql.NullString
			acceptedAt         pq.NullTime
			rejectedAt         pq.NullTime
			createdBy          sql.NullInt64
			updatedBy          sql.NullInt64
			createdByRoleId    sql.NullInt64
			createdByRoleName  sql.NullString
			updatedByRoleId    sql.NullInt64
			updatedByRoleName  sql.NullString
		)
		err := rows.Scan(
			&holiday.ID,
			&holiday.StartDate,
			&holiday.EndDate,
			&holiday.Title,
			&holiday.PaymentType,
			&description,
			&holiday.Status,
			&rejectedReason,
			&acceptedAt,
			&rejectedAt,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
			&deletedAt,
			&createdBy,
			&updatedBy,
			&createdByID,
			&createdByFirstName,
			&createdByLastName,
			&createdByEmail,
			&createdByAvatarUrl,
			&updatedByID,
			&updatedByFirstName,
			&updatedByLastName,
			&updatedByEmail,
			&updatedByAvatarUrl,
			&createdByRoleId,
			&createdByRoleName,
			&updatedByRoleId,
			&updatedByRoleName,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			holiday.Description = &description.String
		}
		if deletedAt.Valid {
			holiday.DeletedAt = &deletedAt.Time
		}
		if rejectedReason.Valid {
			holiday.RejectedReason = &rejectedReason.String
		}
		if acceptedAt.Valid {
			holiday.AcceptedAt = &acceptedAt.Time
		}
		if rejectedAt.Valid {
			holiday.RejectedAt = &rejectedAt.Time
		}
		if createdByID.Valid {
			holiday.CreatedBy = domain.HolidayUser{
				ID: uint(createdByID.Int64),
			}
			if createdByFirstName.Valid {
				holiday.CreatedBy.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				holiday.CreatedBy.LastName = createdByLastName.String
			}
			if createdByEmail.Valid {
				holiday.CreatedBy.Email = createdByEmail.String
			}
			if createdByAvatarUrl.Valid {
				holiday.CreatedBy.AvatarUrl = createdByAvatarUrl.String
			}
		}
		if updatedByID.Valid {
			holiday.UpdatedBy = domain.HolidayUser{
				ID: uint(updatedByID.Int64),
			}
			if updatedByFirstName.Valid {
				holiday.UpdatedBy.FirstName = updatedByFirstName.String
			}
			if updatedByLastName.Valid {
				holiday.UpdatedBy.LastName = updatedByLastName.String
			}
			if updatedByEmail.Valid {
				holiday.UpdatedBy.Email = updatedByEmail.String
			}
			if updatedByAvatarUrl.Valid {
				holiday.UpdatedBy.AvatarUrl = updatedByAvatarUrl.String
			}
		}
		if createdByRoleId.Valid {
			holiday.CreatedBy.Role = &domain.HolidayUserRole{
				ID: uint(createdByRoleId.Int64),
			}
			if createdByRoleName.Valid {
				holiday.CreatedBy.Role.Name = createdByRoleName.String
			}
		}
		if updatedByRoleId.Valid {
			holiday.UpdatedBy.Role = &domain.HolidayUserRole{
				ID: uint(updatedByRoleId.Int64),
			}
			if updatedByRoleName.Valid {
				holiday.UpdatedBy.Role.Name = updatedByRoleName.String
			}
		}
		holidays = append(holidays, &holiday)
	}
	return holidays, nil
}

func (r *HolidayRepositoryPostgresDB) UpdateStatus(payload *models.HolidaysUpdateStatusRequestBody, id int64) (*domain.Holiday, error) {
	// Get current time
	currentTime := time.Now()

	// Get the holiday
	holidays, err := r.Query(&models.HolidaysQueryRequestParams{ID: int(id)})
	if err != nil {
		return nil, err
	}
	if len(holidays) == 0 {
		return nil, errors.New("no rows affected")
	}
	holiday := holidays[0]
	if holiday == nil {
		return nil, errors.New("no rows affected")
	}
	if holiday.Status == payload.Status {
		return holiday, nil
	}

	// Update the holiday
	var (
		updatedId  int
		acceptedAt *time.Time
		rejectedAt *time.Time
	)
	if payload.Status == constants.HOLIDAY_STATUS_ACCEPTED {
		acceptedAt = &currentTime
	}
	if payload.Status == constants.HOLIDAY_STATUS_REJECTED {
		rejectedAt = &currentTime
	}
	err = r.PostgresDB.QueryRow(`
		UPDATE staffClubHolidays
		SET status = $1, accepted_at = $2, rejected_at = $3, updated_at = $4, updatedBy = $5, rejectedReason = $6
		WHERE id = $7
		RETURNING id
	`,
		payload.Status,
		acceptedAt,
		rejectedAt,
		currentTime,
		payload.AuthenticatedUser.ID,
		payload.RejectedReason,
		id,
	).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the holiday
	holidays, err = r.Query(&models.HolidaysQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(holidays) == 0 {
		return nil, errors.New("no rows affected")
	}
	return holidays[0], nil
}
