package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/todo/constants"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/todo/domain"
	"github.com/hoitek/Maja-Service/internal/todo/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type TodoRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewTodoRepositoryPostgresDB(d *sql.DB) *TodoRepositoryPostgresDB {
	return &TodoRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.TodosQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" t.id = %d", queries.ID))
		}
		if queries.UserID != 0 {
			where = append(where, fmt.Sprintf(" t.userId = %d", queries.UserID))
		}
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.Date.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Date.Op, fmt.Sprintf("%v", queries.Filters.Date.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.dateValue %s %s", opValue.Operator, val))
		}
		if queries.Filters.Time.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Time.Op, fmt.Sprintf("%v", queries.Filters.Time.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.timeValue %s %s", opValue.Operator, val))
		}
		if queries.Filters.Description.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Description.Op, fmt.Sprintf("%v", queries.Filters.Description.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.description %s %s", opValue.Operator, val))
		}
		if queries.Filters.Status.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Status.Op, fmt.Sprintf("%v", queries.Filters.Status.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" t.status %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *TodoRepositoryPostgresDB) Query(queries *models.TodosQueryRequestParams) ([]*domain.Todo, error) {
	q := `
		SELECT
			t.id,
			t.userId,
			t.title,
			t.timeValue,
			t.dateValue,
			t.description,
			t.status,
			t.done_at,
			t.created_at,
			t.updated_at,
			t.deleted_at,
			u.id AS userId,
			u.firstName AS userFirstName,
			u.lastName AS userLastName,
			u.avatarUrl AS userAvatarUrl,
			u2.id AS createdById,
			u2.firstName AS createdByFirstName,
			u2.lastName AS createdByLastName,
			u2.avatarUrl AS createdByAvatarUrl
		FROM todos t
		LEFT JOIN users u ON t.userId = u.id
		LEFT JOIN users u2 ON t.createdBy = u2.id
	`
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.name %s", queries.Sorts.Title.Op))
		}
		if queries.Sorts.Date.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.date %s", queries.Sorts.Date.Op))
		}
		if queries.Sorts.Time.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.time %s", queries.Sorts.Time.Op))
		}
		if queries.Sorts.UserID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.userId %s", queries.Sorts.UserID.Op))
		}
		if queries.Sorts.Description.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.description %s", queries.Sorts.Description.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" t.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var todos []*domain.Todo
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			todo               domain.Todo
			dontAt             sql.NullTime
			deletedAt          sql.NullTime
			userId             sql.NullInt64
			userFirstName      sql.NullString
			userLastName       sql.NullString
			userAvatarUrl      sql.NullString
			createdById        sql.NullInt64
			createdByFirstName sql.NullString
			createdByLastName  sql.NullString
			createdByAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Time,
			&todo.Date,
			&todo.Description,
			&todo.Status,
			&dontAt,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&createdById,
			&createdByFirstName,
			&createdByLastName,
			&createdByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		todo.DoneAt = exp.TerIf(dontAt.Valid, &dontAt.Time, nil)
		todo.DeletedAt = exp.TerIf(deletedAt.Valid, &deletedAt.Time, nil)
		if userId.Valid {
			todo.User = &domain.TodoUser{
				ID: uint(userId.Int64),
			}
			if userFirstName.Valid {
				todo.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				todo.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				todo.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if createdById.Valid {
			todo.CreatedBy = uint(createdById.Int64)
			todo.CreatedByUser = &domain.TodoUser{
				ID: todo.CreatedBy,
			}
			if createdByFirstName.Valid {
				todo.CreatedByUser.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				todo.CreatedByUser.LastName = createdByLastName.String
			}
			if createdByAvatarUrl.Valid {
				todo.CreatedByUser.AvatarUrl = createdByAvatarUrl.String
			}
		}

		// Format Date
		todo.DateStr = todo.Date.Format("2006-01-02")

		// Format Time
		todo.TimeStr = todo.Time.Format("15:04:05")

		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *TodoRepositoryPostgresDB) Count(queries *models.TodosQueryRequestParams) (int64, error) {
	q := `
		SELECT
		    COUNT(t.id)
		FROM todos t
		LEFT JOIN users u ON t.userId = u.id
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

func (r *TodoRepositoryPostgresDB) Create(payload *models.TodosCreateRequestBody) (*domain.Todo, error) {
	var (
		insertedID int
		status     = constants.TODO_STATUS_ACTIVE
	)

	// Current time
	currentTime := time.Now()

	// Insert the todos
	err := r.PostgresDB.QueryRow(`
		INSERT INTO todos (userId, title, timeValue, dateValue, description, status, done_at, created_at, updated_at, deleted_at, createdBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		payload.UserID,
		payload.Title,
		payload.Time,
		payload.Date,
		payload.Description,
		status,
		nil,
		currentTime,
		currentTime,
		nil,
		payload.AuthenticatedUser.ID,
	).Scan(
		&insertedID,
	)
	if err != nil {
		return nil, err
	}

	// Find the todos by id
	todos, err := r.Query(&models.TodosQueryRequestParams{
		ID: insertedID,
	})
	if err != nil {
		return nil, err
	}
	if len(todos) == 0 {
		return nil, errors.New("no todos found")
	}
	todo := todos[0]

	// Return the todos
	return todo, nil
}

func (r *TodoRepositoryPostgresDB) Delete(payload *models.TodosDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM todos
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

func (r *TodoRepositoryPostgresDB) Update(payload *models.TodosCreateRequestBody, id int64) (*domain.Todo, error) {
	var (
		updatedID int
		status    = constants.TODO_STATUS_ACTIVE
	)

	// Current time
	currentTime := time.Now()

	// Find the todos by id
	foundTodos, err := r.Query(&models.TodosQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(foundTodos) == 0 {
		return nil, errors.New("no todos found")
	}
	foundTodo := foundTodos[0]

	// Update the todos
	err = r.PostgresDB.QueryRow(`
		UPDATE todos
		SET userId = $1, title = $2, timeValue = $3, dateValue = $4, description = $5, status = $6, updated_at = $7
		WHERE id = $8
		RETURNING id
	`,
		payload.UserID,
		payload.Title,
		payload.Time,
		payload.Date,
		payload.Description,
		status,
		currentTime,
		foundTodo.ID,
	).Scan(&updatedID)

	// If the todos does not update, return an error
	if err != nil {
		return nil, err
	}

	// Find the todos by id
	todos, err := r.Query(&models.TodosQueryRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}
	if len(todos) == 0 {
		return nil, errors.New("no todos found")
	}
	todo := todos[0]

	// Return the todos
	return todo, nil
}

func (r *TodoRepositoryPostgresDB) UpdateStatus(payload *models.TodosUpdateStatusRequestBody, id int64) (*domain.Todo, error) {
	var (
		updatedID int
		doneAt    *time.Time = nil
	)

	// Current time
	currentTime := time.Now()

	// Find the todos by id
	foundTodos, err := r.Query(&models.TodosQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(foundTodos) == 0 {
		return nil, errors.New("no todos found")
	}
	foundTodo := foundTodos[0]

	if payload.Status == constants.TODO_STATUS_DONE {
		doneAt = &currentTime
	}

	// Update the todos
	err = r.PostgresDB.QueryRow(`
		UPDATE todos
		SET status = $1, done_at = $2, updated_at = $3
		WHERE id = $4
		RETURNING id
	`,
		payload.Status,
		doneAt,
		currentTime,
		foundTodo.ID,
	).Scan(&updatedID)

	// If the todos does not update, return an error
	if err != nil {
		return nil, err
	}

	// Find the todos by id
	todos, err := r.Query(&models.TodosQueryRequestParams{
		ID: updatedID,
	})
	if err != nil {
		return nil, err
	}
	if len(todos) == 0 {
		return nil, errors.New("no todos found")
	}
	todo := todos[0]

	// Return the todos
	return todo, nil
}

func (r *TodoRepositoryPostgresDB) GetTodosByIds(ids []int64) ([]*domain.Todo, error) {
	var todos []*domain.Todo
	rows, err := r.PostgresDB.Query(`
		SELECT
			t.id,
			t.userId,
			t.title,
			t.timeValue,
			t.dateValue,
			t.description,
			t.status,
			t.done_at,
			t.created_at,
			t.updated_at,
			t.deleted_at,
			u.id AS userId,
			u.firstName AS userFirstName,
			u.lastName AS userLastName,
			u.avatarUrl AS userAvatarUrl,
			u2.id AS createdById,
			u2.firstName AS createdByFirstName,
			u2.lastName AS createdByLastName,
			u2.avatarUrl AS createdByAvatarUrl
		FROM todos t
		LEFT JOIN users u ON t.userId = u.id
		LEFT JOIN users u2 ON t.createdBy = u2.id
		WHERE t.id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			todo               domain.Todo
			dontAt             sql.NullTime
			deletedAt          sql.NullTime
			userId             sql.NullInt64
			userFirstName      sql.NullString
			userLastName       sql.NullString
			userAvatarUrl      sql.NullString
			createdById        sql.NullInt64
			createdByFirstName sql.NullString
			createdByLastName  sql.NullString
			createdByAvatarUrl sql.NullString
		)
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Time,
			&todo.Date,
			&todo.Description,
			&todo.Status,
			&dontAt,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&deletedAt,
			&userId,
			&userFirstName,
			&userLastName,
			&userAvatarUrl,
			&createdById,
			&createdByFirstName,
			&createdByLastName,
			&createdByAvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		todo.DoneAt = exp.TerIf(dontAt.Valid, &dontAt.Time, nil)
		todo.DeletedAt = exp.TerIf(deletedAt.Valid, &deletedAt.Time, nil)
		if userId.Valid {
			todo.User = &domain.TodoUser{
				ID: uint(userId.Int64),
			}
			if userFirstName.Valid {
				todo.User.FirstName = userFirstName.String
			}
			if userLastName.Valid {
				todo.User.LastName = userLastName.String
			}
			if userAvatarUrl.Valid {
				todo.User.AvatarUrl = userAvatarUrl.String
			}
		}
		if createdById.Valid {
			todo.CreatedBy = uint(createdById.Int64)
			todo.CreatedByUser = &domain.TodoUser{
				ID: todo.CreatedBy,
			}
			if createdByFirstName.Valid {
				todo.CreatedByUser.FirstName = createdByFirstName.String
			}
			if createdByLastName.Valid {
				todo.CreatedByUser.LastName = createdByLastName.String
			}
			if createdByAvatarUrl.Valid {
				todo.CreatedByUser.AvatarUrl = createdByAvatarUrl.String
			}
		}

		// Format Date
		todo.DateStr = todo.Date.Format("2006-01-02")

		// Format Time
		todo.TimeStr = todo.Time.Format("15:04:05")

		todos = append(todos, &todo)
	}
	return todos, nil
}
