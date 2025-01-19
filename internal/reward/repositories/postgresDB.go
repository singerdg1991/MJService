package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/reward/domain"
	"github.com/hoitek/Maja-Service/internal/reward/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type RewardRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewRewardRepositoryPostgresDB(d *sql.DB) *RewardRepositoryPostgresDB {
	return &RewardRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.RewardsQueryRequestParams) []string {
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

func (r *RewardRepositoryPostgresDB) Query(queries *models.RewardsQueryRequestParams) ([]*domain.Reward, error) {
	q := `SELECT * FROM rewards ls `
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

	var rewards []*domain.Reward
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reward domain.Reward
		err := rows.Scan(&reward.ID, &reward.Name, &reward.Description, &reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, &reward)
	}
	return rewards, nil
}

func (r *RewardRepositoryPostgresDB) Count(queries *models.RewardsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM rewards ls `
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

func (r *RewardRepositoryPostgresDB) Create(payload *models.RewardsCreateRequestBody) (*domain.Reward, error) {
	var reward domain.Reward

	// Current time
	currentTime := time.Now()

	// Insert the reward
	err := r.PostgresDB.QueryRow(`
		INSERT INTO rewards (name, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, payload.Description, currentTime, currentTime, nil).Scan(&reward.ID, &reward.Name, &reward.Description, &reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the reward
	return &reward, nil
}

func (r *RewardRepositoryPostgresDB) Delete(payload *models.RewardsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM rewards
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

func (r *RewardRepositoryPostgresDB) Update(payload *models.RewardsCreateRequestBody, id int64) (*domain.Reward, error) {
	var reward domain.Reward

	// Current time
	currentTime := time.Now()

	// Find the reward by name
	var foundReward domain.Reward
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM rewards
		WHERE id = $1
	`, id).Scan(&foundReward.ID, &foundReward.Name, &foundReward.Description, &foundReward.CreatedAt, &foundReward.UpdatedAt, &foundReward.DeletedAt)

	// If the reward is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the reward
	err = r.PostgresDB.QueryRow(`
		UPDATE rewards
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, payload.Description, foundReward.ID).Scan(&reward.ID, &reward.Name, &reward.Description, &reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt)

	// If the reward does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the reward
	return &reward, nil
}

func (r *RewardRepositoryPostgresDB) GetRewardsByIds(ids []int64) ([]*domain.Reward, error) {
	var rewards []*domain.Reward
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM rewards
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reward domain.Reward
		err := rows.Scan(&reward.ID, &reward.Name, &reward.Description, &reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, &reward)
	}
	return rewards, nil
}
