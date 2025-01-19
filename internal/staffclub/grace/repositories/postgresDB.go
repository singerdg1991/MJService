package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type GraceRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewGraceRepositoryPostgresDB(d *sql.DB) *GraceRepositoryPostgresDB {
	return &GraceRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.GracesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" g.id = %d", queries.ID))
		}
		if queries.Filters.GraceNumber.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.GraceNumber.Op, fmt.Sprintf("%v", queries.Filters.GraceNumber.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.graceNumber %s %s", opValue.Operator, val))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" g.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *GraceRepositoryPostgresDB) Query(queries *models.GracesQueryRequestParams) ([]*domain.Grace, error) {
	q := `
		SELECT
		    g.id,
		    g.rewardId,
			g.isAutoRewardSetEnable,
			g.graceNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS rewardId,
			r.name AS rewardName,
			r.description AS rewardDescription
		FROM staffClubGraces g
		LEFT JOIN rewards r ON g.rewardId = r.id
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

	var graces []*domain.Grace
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			grace             domain.Grace
			description       sql.NullString
			deletedAt         pq.NullTime
			rewardId          sql.NullInt64
			rewardName        sql.NullString
			rewardDescription sql.NullString
		)
		err := rows.Scan(
			&grace.ID,
			&grace.RewardID,
			&grace.IsAutoRewardSetEnable,
			&grace.GraceNumber,
			&grace.Title,
			&description,
			&grace.CreatedAt,
			&grace.UpdatedAt,
			&deletedAt,
			&rewardId,
			&rewardName,
			&rewardDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			grace.Description = &description.String
		}
		if deletedAt.Valid {
			grace.DeletedAt = &deletedAt.Time
		}
		if rewardId.Valid {
			grace.Reward = &domain.GraceReward{
				ID: uint(rewardId.Int64),
			}
			if rewardName.Valid {
				grace.Reward.Name = rewardName.String
			}
			if rewardDescription.Valid {
				grace.Reward.Description = rewardDescription.String
			}
		}
		graces = append(graces, &grace)
	}
	return graces, nil
}

func (r *GraceRepositoryPostgresDB) Count(queries *models.GracesQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(g.id)
		FROM staffClubGraces g
		LEFT JOIN rewards r ON g.rewardId = r.id
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

func (r *GraceRepositoryPostgresDB) Create(payload *models.GracesCreateRequestBody) (*domain.Grace, error) {
	// Current time
	currentTime := time.Now()

	// Insert the grace
	var insertedId int
	err := r.PostgresDB.QueryRow(`
		INSERT INTO staffClubGraces (rewardId, isAutoRewardSetEnable, graceNumber, title, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.RewardID,
		payload.IsAutoRewardSetEnableAsBool,
		payload.GraceNumber,
		payload.Title,
		payload.Description,
		currentTime,
		currentTime,
		nil,
	).Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	// Get the grace
	graces, err := r.Query(&models.GracesQueryRequestParams{ID: insertedId})
	if err != nil {
		return nil, err
	}
	if len(graces) == 0 {
		return nil, errors.New("no rows affected")
	}
	return graces[0], nil
}

func (r *GraceRepositoryPostgresDB) Delete(payload *models.GracesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM staffClubGraces
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

func (r *GraceRepositoryPostgresDB) Update(payload *models.GracesCreateRequestBody, id int64) (*domain.Grace, error) {
	// Current time
	currentTime := time.Now()

	// Find the grace by id
	foundGraces, _ := r.Query(&models.GracesQueryRequestParams{ID: int(id)})
	if len(foundGraces) == 0 {
		return nil, errors.New("grace not found")
	}
	foundGrace := foundGraces[0]
	if foundGrace.DeletedAt != nil {
		return nil, errors.New("grace is deleted")
	}

	// Update the grace
	var updatedId int
	err := r.PostgresDB.QueryRow(`
		UPDATE staffClubGraces
		SET rewardId = $1, isAutoRewardSetEnable = $2, graceNumber = $3, title = $4, description = $5, updated_at = $6
		WHERE id = $7
		RETURNING id
	`, payload.RewardID, payload.IsAutoRewardSetEnableAsBool, payload.GraceNumber, payload.Title, payload.Description, currentTime, id).Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	// Get the grace
	graces, err := r.Query(&models.GracesQueryRequestParams{ID: updatedId})
	if err != nil {
		return nil, err
	}
	if len(graces) == 0 {
		return nil, errors.New("no rows affected")
	}
	return graces[0], nil
}

func (r *GraceRepositoryPostgresDB) GetGracesByIds(ids []int64) ([]*domain.Grace, error) {
	var graces []*domain.Grace
	rows, err := r.PostgresDB.Query(`
		SELECT
		    g.id,
		    g.rewardId,
			g.isAutoRewardSetEnable,
			g.graceNumber,
			g.title,
			g.description,
			g.created_at,
			g.updated_at,
			g.deleted_at,
			r.id AS rewardId,
			r.name AS rewardName,
			r.description AS rewardDescription
		FROM staffClubGraces g
		LEFT JOIN rewards r ON g.rewardId = r.id
		WHERE g.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			grace             domain.Grace
			description       sql.NullString
			deletedAt         pq.NullTime
			rewardId          sql.NullInt64
			rewardName        sql.NullString
			rewardDescription sql.NullString
		)
		err := rows.Scan(
			&grace.ID,
			&grace.RewardID,
			&grace.IsAutoRewardSetEnable,
			&grace.GraceNumber,
			&grace.Title,
			&description,
			&grace.CreatedAt,
			&grace.UpdatedAt,
			&deletedAt,
			&rewardId,
			&rewardName,
			&rewardDescription,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			grace.Description = &description.String
		}
		if deletedAt.Valid {
			grace.DeletedAt = &deletedAt.Time
		}
		if rewardId.Valid {
			grace.Reward = &domain.GraceReward{
				ID: uint(rewardId.Int64),
			}
			if rewardName.Valid {
				grace.Reward.Name = rewardName.String
			}
			if rewardDescription.Valid {
				grace.Reward.Description = rewardDescription.String
			}
		}
		graces = append(graces, &grace)
	}
	return graces, nil
}
