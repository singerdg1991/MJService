package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/languageskill/domain"
	"github.com/hoitek/Maja-Service/internal/languageskill/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type LanguageSkillRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewLanguageSkillRepositoryPostgresDB(d *sql.DB) *LanguageSkillRepositoryPostgresDB {
	return &LanguageSkillRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.LanguageSkillsQueryRequestParams) []string {
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

func (r *LanguageSkillRepositoryPostgresDB) Query(queries *models.LanguageSkillsQueryRequestParams) ([]*domain.LanguageSkill, error) {
	q := `SELECT * FROM languageskills ls `
	log.Println(queries)
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
	log.Println(q)

	var languageskills []*domain.LanguageSkill
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			languageskill domain.LanguageSkill
			description   sql.NullString
			deletedAt     pq.NullTime
		)
		err := rows.Scan(&languageskill.ID, &languageskill.Name, &languageskill.CreatedAt, &languageskill.UpdatedAt, &deletedAt, &description)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			languageskill.DeletedAt = &deletedAt.Time
		}
		if description.Valid {
			languageskill.Description = &description.String
		}
		languageskills = append(languageskills, &languageskill)
	}
	return languageskills, nil
}

func (r *LanguageSkillRepositoryPostgresDB) Count(queries *models.LanguageSkillsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM languageskills ls `
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

func (r *LanguageSkillRepositoryPostgresDB) Create(payload *models.LanguageSkillsCreateRequestBody) (*domain.LanguageSkill, error) {
	var languageskill domain.LanguageSkill

	// Current time
	currentTime := time.Now()

	// Insert the languageskill
	err := r.PostgresDB.QueryRow(`
		INSERT INTO languageskills (name, created_at, updated_at, deleted_at, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, created_at, updated_at, deleted_at, description
	`, payload.Name, currentTime, currentTime, nil, payload.Description).Scan(&languageskill.ID, &languageskill.Name, &languageskill.CreatedAt, &languageskill.UpdatedAt, &languageskill.DeletedAt, &languageskill.Description)
	if err != nil {
		return nil, err
	}

	// Return the languageskill
	return &languageskill, nil
}

func (r *LanguageSkillRepositoryPostgresDB) Delete(payload *models.LanguageSkillsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM languageskills
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

func (r *LanguageSkillRepositoryPostgresDB) Update(payload *models.LanguageSkillsCreateRequestBody, id int64) (*domain.LanguageSkill, error) {
	var languageskill domain.LanguageSkill

	// Current time
	currentTime := time.Now()

	// Find the languageskill by name
	var foundLanguageSkill domain.LanguageSkill
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM languageskills
		WHERE id = $1
	`, id).Scan(&foundLanguageSkill.ID, &foundLanguageSkill.Name, &foundLanguageSkill.CreatedAt, &foundLanguageSkill.UpdatedAt, &foundLanguageSkill.DeletedAt, &foundLanguageSkill.Description)

	// If the languageskill is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the languageskill
	err = r.PostgresDB.QueryRow(`
		UPDATE languageskills
		SET name = $1, updated_at = $2, description = $3
		WHERE id = $4
		RETURNING id, name, created_at, updated_at, deleted_at, description
	`, payload.Name, currentTime, payload.Description, foundLanguageSkill.ID).Scan(&languageskill.ID, &languageskill.Name, &languageskill.CreatedAt, &languageskill.UpdatedAt, &languageskill.DeletedAt, &languageskill.Description)

	// If the languageskill does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the languageskill
	return &languageskill, nil
}

func (r *LanguageSkillRepositoryPostgresDB) GetLanguageSkillsByIds(ids []int64) ([]*domain.LanguageSkill, error) {
	var languageSkills []*domain.LanguageSkill
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM languageSkills
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			languageSkill domain.LanguageSkill
			description   sql.NullString
			deletedAt     pq.NullTime
		)
		err := rows.Scan(&languageSkill.ID, &languageSkill.Name, &languageSkill.CreatedAt, &languageSkill.UpdatedAt, &deletedAt, &description)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			languageSkill.DeletedAt = &deletedAt.Time
		}
		if description.Valid {
			languageSkill.Description = &description.String
		}
		languageSkills = append(languageSkills, &languageSkill)
	}
	return languageSkills, nil
}
