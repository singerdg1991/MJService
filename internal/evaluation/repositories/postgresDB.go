package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/evaluation/domain"
	"github.com/hoitek/Maja-Service/internal/evaluation/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type EvaluationRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewEvaluationRepositoryPostgresDB(d *sql.DB) *EvaluationRepositoryPostgresDB {
	return &EvaluationRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.EvaluationsQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" ls.id = %d", queries.ID))
		}
		if queries.Filters.StaffID.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.StaffID.Op, fmt.Sprintf("%v", queries.Filters.StaffID.Value))
			val := exp.TerIf(opValue.Value == "", "", opValue.Value)
			where = append(where, fmt.Sprintf(" ls.staffId %s %s", opValue.Operator, val))
		}
		if queries.Filters.EvaluationType.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.EvaluationType.Op, fmt.Sprintf("%v", queries.Filters.EvaluationType.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.evaluationType %s %s", opValue.Operator, val))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.title %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.description %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *EvaluationRepositoryPostgresDB) Query(queries *models.EvaluationsQueryRequestParams) ([]*domain.Evaluation, error) {
	q := `SELECT * FROM evaluations ls `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.EvaluationType.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.evaluationType %s", queries.Sorts.EvaluationType.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.title %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.description %s", queries.Sorts.Description.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var evaluations []*domain.Evaluation
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			evaluation  domain.Evaluation
			description sql.NullString
			deletedAt   pq.NullTime
		)
		err := rows.Scan(
			&evaluation.ID,
			&evaluation.StaffID,
			&evaluation.EvaluationType,
			&evaluation.Title,
			&description,
			&evaluation.CreatedAt,
			&evaluation.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			evaluation.Description = &description.String
		}
		if deletedAt.Valid {
			evaluation.DeletedAt = &deletedAt.Time
		}
		evaluations = append(evaluations, &evaluation)
	}
	return evaluations, nil
}

func (r *EvaluationRepositoryPostgresDB) Count(queries *models.EvaluationsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM evaluations ls `
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

func (r *EvaluationRepositoryPostgresDB) Create(payload *models.EvaluationsCreateRequestBody) (*domain.Evaluation, error) {
	var evaluation domain.Evaluation

	// Current time
	currentTime := time.Now()

	// Insert the evaluation
	var (
		description sql.NullString
		deletedAt   sql.NullTime
	)
	err := r.PostgresDB.QueryRow(`
		INSERT INTO evaluations (staffId, evaluationType, title, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, staffId, evaluationType, title, description, created_at, updated_at, deleted_at
	`, payload.StaffID, payload.EvaluationType, payload.Title, payload.Description, currentTime, currentTime, nil).Scan(
		&evaluation.ID,
		&evaluation.StaffID,
		&evaluation.EvaluationType,
		&evaluation.Title,
		&description,
		&evaluation.CreatedAt,
		&evaluation.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		evaluation.Description = &description.String
	}
	if deletedAt.Valid {
		evaluation.DeletedAt = &deletedAt.Time
	}

	// Return the evaluation
	return &evaluation, nil
}

func (r *EvaluationRepositoryPostgresDB) Delete(payload *models.EvaluationsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM evaluations
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

func (r *EvaluationRepositoryPostgresDB) Update(payload *models.EvaluationsCreateRequestBody, id int64) (*domain.Evaluation, error) {
	var evaluation domain.Evaluation

	// Current time
	currentTime := time.Now()

	// Update the evaluation
	var (
		description sql.NullString
		deletedAt   sql.NullTime
	)
	err := r.PostgresDB.QueryRow(`
		UPDATE evaluations
		SET staffId = $1, evaluationType = $2, title = $3, updated_at = $4, description = $5
		WHERE id = $6
		RETURNING id, staffId, evaluationType, title, description, created_at, updated_at, deleted_at
	`, payload.StaffID, payload.EvaluationType, payload.Title, currentTime, payload.Description, id).Scan(
		&evaluation.ID,
		&evaluation.StaffID,
		&evaluation.EvaluationType,
		&evaluation.Title,
		&description,
		&evaluation.CreatedAt,
		&evaluation.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		evaluation.Description = &description.String
	}
	if deletedAt.Valid {
		evaluation.DeletedAt = &deletedAt.Time
	}

	// Return the evaluation
	return &evaluation, nil
}

func (r *EvaluationRepositoryPostgresDB) GetEvaluationsByIds(ids []int64) ([]*domain.Evaluation, error) {
	var evaluations []*domain.Evaluation
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM evaluations
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			evaluation  domain.Evaluation
			description sql.NullString
			deletedAt   pq.NullTime
		)
		err := rows.Scan(
			&evaluation.ID,
			&evaluation.StaffID,
			&evaluation.EvaluationType,
			&evaluation.Title,
			&description,
			&evaluation.CreatedAt,
			&evaluation.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			evaluation.Description = &description.String
		}
		if deletedAt.Valid {
			evaluation.DeletedAt = &deletedAt.Time
		}
		evaluations = append(evaluations, &evaluation)
	}
	return evaluations, nil
}
