package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	quilder "github.com/hoitek/Go-Quilder"
	"github.com/hoitek/Maja-Service/internal/ability/domain"
	"github.com/hoitek/Maja-Service/internal/ability/models"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

type AbilityRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewAbilityRepositoryPostgresDB(d *sql.DB) *AbilityRepositoryPostgresDB {
	return &AbilityRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func (r *AbilityRepositoryPostgresDB) Query(queries *models.AbilitiesQueryRequestParams) ([]*domain.Ability, error) {
	q := `SELECT * FROM abilities `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit

		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return nil, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var abilities []*domain.Ability
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ability domain.Ability
		err := rows.Scan(&ability.ID, &ability.Name, &ability.CreatedAt, &ability.UpdatedAt, &ability.DeletedAt)
		if err != nil {
			return nil, err
		}
		abilities = append(abilities, &ability)
	}
	return abilities, nil
}

func (r *AbilityRepositoryPostgresDB) Count(queries *models.AbilitiesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM abilities `
	if queries != nil {
		limit := queries.Limit
		offset := (queries.Page - 1) * limit
		options, err := quilder.CreateQueriesGroup(queries)
		if err != nil {
			return 0, err
		}
		qo := &quilder.Query{
			QueriesGroup: options,
			Limit:        queries.Limit,
			Offset:       offset,
		}
		query := qo.Build()
		if query.Where != "" || queries.ID > 0 {
			q += "WHERE "
		}
		if queries.ID > 0 {
			q += fmt.Sprintf("id = %d", queries.ID)
			if query.Where != "" {
				q += " AND "
			}
		}
		q += query.Where
		if query.Limit != 0 {
			q += fmt.Sprintf(" LIMIT %d", query.Limit)
		}
		if query.Offset != 0 {
			q += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	var count int64
	err := r.PostgresDB.QueryRow(q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AbilityRepositoryPostgresDB) Create(payload *models.AbilitiesCreateRequestBody) (*domain.Ability, error) {
	var ability domain.Ability

	// Current time
	currentTime := time.Now()

	// Find the ability by name
	var foundAbility domain.Ability
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM abilities
WHERE name = $1
`, payload.Name).Scan(&foundAbility.ID, &foundAbility.Name, &foundAbility.CreatedAt, &foundAbility.UpdatedAt, &foundAbility.DeletedAt)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the ability
	err = r.PostgresDB.QueryRow(`
INSERT INTO abilities (name, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4)
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, currentTime, nil).Scan(&ability.ID, &ability.Name, &ability.CreatedAt, &ability.UpdatedAt, &ability.DeletedAt)
	// If the ability does not insert, return an error
	if err != nil {
		return nil, err
	}

	// Return the ability
	return &ability, nil
}

func (r *AbilityRepositoryPostgresDB) Delete(payload *models.AbilitiesDeleteRequestBody) ([]int64, error) {
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
DELETE FROM abilities
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

func (r *AbilityRepositoryPostgresDB) Update(payload *models.AbilitiesCreateRequestBody, name string) (*domain.Ability, error) {
	var ability domain.Ability

	// Current time
	currentTime := time.Now()

	// Find the ability by name
	var foundAbility domain.Ability
	err := r.PostgresDB.QueryRow(`
SELECT *
FROM abilities
WHERE name = $1
`, name).Scan(&foundAbility.ID, &foundAbility.Name, &foundAbility.CreatedAt, &foundAbility.UpdatedAt, &foundAbility.DeletedAt)

	// If the ability is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the ability
	err = r.PostgresDB.QueryRow(`
UPDATE abilities
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING id, name, created_at, updated_at, deleted_at
`, payload.Name, currentTime, foundAbility.ID).Scan(&ability.ID, &ability.Name, &ability.CreatedAt, &ability.UpdatedAt, &ability.DeletedAt)

	// If the ability does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the ability
	return &ability, nil
}

func (r *AbilityRepositoryPostgresDB) GetAbilitiesByIds(ids []int64) ([]*domain.Ability, error) {
	var abilities []*domain.Ability
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM abilities
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ability domain.Ability
		err := rows.Scan(&ability.ID, &ability.Name, &ability.CreatedAt, &ability.UpdatedAt, &ability.DeletedAt)
		if err != nil {
			return nil, err
		}
		abilities = append(abilities, &ability)
	}
	return abilities, nil
}
