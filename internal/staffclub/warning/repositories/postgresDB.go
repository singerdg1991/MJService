package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type WarningRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewWarningRepositoryPostgresDB(d *sql.DB) *WarningRepositoryPostgresDB {
	return &WarningRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.WarningsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" g.id = %d", queries.ID))
		}
		if queries.Filters.WarningNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.WarningNumber.Op, fmt.Sprintf("%v", queries.Filters.WarningNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.warningNumber %s %s", opValue.Operator, val))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *WarningRepositoryPostgresDB) Query(queries *models.WarningsQueryRequestParams) ([]*domain.Warning, error) {
	q := `
		SELECT
		    g.id,
		    g.punishmentId,
			g.isAutoRewardSetEnable,
			g.warningNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS punishmentId,
			r.name AS punishmentName,
			r.description AS punishmentDescription
		FROM staffClubWarnings g
		LEFT JOIN punishments r ON g.punishmentId = r.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" g.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var warnings []*domain.Warning
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			warning               domain.Warning
			description           sql.NullString
			deletedAt             pq.NullTime
			punishmentId          sql.NullInt64
			punishmentName        sql.NullString
			punishmentDescription sql.NullString
		)
		err := rows.Scan(
			&warning.ID,
			&warning.PunishmentID,
			&warning.IsAutoRewardSetEnable,
			&warning.WarningNumber,
			&warning.Title,
			&description,
			&warning.CreatedAt,
			&warning.UpdatedAt,
			&deletedAt,
			&punishmentId,
			&punishmentName,
			&punishmentDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			warning.Description = &description.String
		}
		if deletedAt.Valid {
			warning.DeletedAt = &deletedAt.Time
		}
		if punishmentId.Valid {
			warning.Punishment = &domain.WarningPunishment{
				ID: uint(punishmentId.Int64),
			}
			if punishmentName.Valid {
				warning.Punishment.Name = punishmentName.String
			}
			if punishmentDescription.Valid {
				warning.Punishment.Description = punishmentDescription.String
			}
		}
		warnings = append(warnings, &warning)
	}
	return warnings, nil
}

func (r *WarningRepositoryPostgresDB) Count(queries *models.WarningsQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(g.id)
		FROM staffClubWarnings g
		LEFT JOIN punishments r ON g.punishmentId = r.id
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

func (r *WarningRepositoryPostgresDB) Create(payload *models.WarningsCreateRequestBody) (*domain.Warning, error) {
	// Current time
	currentTime := time.Now()

	// Insert the warning
	var insertedId int
	err := r.PostgresDB.QueryRow(`
		INSERT INTO staffClubWarnings (punishmentId, isAutoRewardSetEnable, warningNumber, title, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.PunishmentID,
		payload.IsAutoRewardSetEnableAsBool,
		payload.WarningNumber,
		payload.Title,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the warning
	warnings, err := r.Query(&models.WarningsQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(warnings) == 0 {
		return nil, errors.New("no rows affected")
	}
	return warnings[0], nil
}

func (r *WarningRepositoryPostgresDB) Delete(payload *models.WarningsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM staffClubWarnings
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

func (r *WarningRepositoryPostgresDB) Update(payload *models.WarningsCreateRequestBody, id int64) (*domain.Warning, error) {
	// Current time
	currentTime := time.Now()

	// Find the warning by id
	foundWarnings, _ := r.Query(&models.WarningsQueryRequestParams{ID: int(id)})
	if len(foundWarnings) == 0 {
		return nil, errors.New("warning not found")
	}
	foundWarning := foundWarnings[0]
	if foundWarning.DeletedAt != nil {
		return nil, errors.New("warning is deleted")
	}

	// Update the warning
	var updatedId int
	err := r.PostgresDB.QueryRow(`
		UPDATE staffClubWarnings
		SET punishmentId = $1, isAutoRewardSetEnable = $2, warningNumber = $3, title = $4, description = $5, updated_at = $6
		WHERE id = $7
		RETURNING id
	`, payload.PunishmentID, payload.IsAutoRewardSetEnableAsBool, payload.WarningNumber, payload.Title, payload.Description, currentTime, id).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the warning
	warnings, err := r.Query(&models.WarningsQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(warnings) == 0 {
		return nil, errors.New("no rows affected")
	}
	return warnings[0], nil
}

func (r *WarningRepositoryPostgresDB) GetWarningsByIds(ids []int64) ([]*domain.Warning, error) {
	var warnings []*domain.Warning
	rows, err := r.PostgresDB.Query(`
		SELECT
		    g.id,
		    g.punishmentId,
			g.isAutoRewardSetEnable,
			g.warningNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS punishmentId,
			r.name AS punishmentName,
			r.description AS punishmentDescription
		FROM staffClubWarnings g
		LEFT JOIN punishments r ON g.punishmentId = r.id
		WHERE g.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			warning               domain.Warning
			description           sql.NullString
			deletedAt             pq.NullTime
			punishmentId          sql.NullInt64
			punishmentName        sql.NullString
			punishmentDescription sql.NullString
		)
		err := rows.Scan(
			&warning.ID,
			&warning.PunishmentID,
			&warning.IsAutoRewardSetEnable,
			&warning.WarningNumber,
			&warning.Title,
			&description,
			&warning.CreatedAt,
			&warning.UpdatedAt,
			&deletedAt,
			&punishmentId,
			&punishmentName,
			&punishmentDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			warning.Description = &description.String
		}
		if deletedAt.Valid {
			warning.DeletedAt = &deletedAt.Time
		}
		if punishmentId.Valid {
			warning.Punishment = &domain.WarningPunishment{
				ID: uint(punishmentId.Int64),
			}
			if punishmentName.Valid {
				warning.Punishment.Name = punishmentName.String
			}
			if punishmentDescription.Valid {
				warning.Punishment.Description = punishmentDescription.String
			}
		}
		warnings = append(warnings, &warning)
	}
	return warnings, nil
}
