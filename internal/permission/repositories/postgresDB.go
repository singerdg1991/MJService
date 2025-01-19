package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/permission/domain"
	"github.com/hoitek/Maja-Service/internal/permission/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type PermissionRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewPermissionRepositoryPostgresDB(d *sql.DB) *PermissionRepositoryPostgresDB {
	return &PermissionRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.PermissionsQueryRequestParams) []string {
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
		if queries.Filters.Title.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Title.Op, fmt.Sprintf("%v", queries.Filters.Title.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" ls.title %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *PermissionRepositoryPostgresDB) Query(queries *models.PermissionsQueryRequestParams) ([]*domain.Permission, error) {
	q := `SELECT * FROM _permissions ls `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.ID.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.id %s", queries.Sorts.ID.Op))
		}
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Title.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" ls.title %s", queries.Sorts.Title.Op))
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

	var permissions []*domain.Permission
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Title, &permission.CreatedAt, &permission.UpdatedAt, &permission.DeletedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

func (r *PermissionRepositoryPostgresDB) Count(queries *models.PermissionsQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(ls.id) FROM _permissions ls `
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

func (r *PermissionRepositoryPostgresDB) Create(payload *models.PermissionsCreateRequestBody) (*domain.Permission, error) {
	var permission domain.Permission

	// Current time
	currentTime := time.Now()

	// Insert the permission
	err := r.PostgresDB.QueryRow(`
		INSERT INTO _permissions (name, title, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, title, created_at, updated_at, deleted_at
	`, payload.Name, payload.Title, currentTime, currentTime, nil).Scan(&permission.ID, &permission.Name, &permission.Title, &permission.CreatedAt, &permission.UpdatedAt, &permission.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Return the permission
	return &permission, nil
}

func (r *PermissionRepositoryPostgresDB) Delete(payload *models.PermissionsDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM _permissions
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

func (r *PermissionRepositoryPostgresDB) Update(payload *models.PermissionsCreateRequestBody, id int64) (*domain.Permission, error) {
	var permission domain.Permission

	// Current time
	currentTime := time.Now()

	// Find the permission by name
	var foundPermission domain.Permission
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM _permissions
		WHERE id = $1
	`, id).Scan(&foundPermission.ID, &foundPermission.Name, &foundPermission.Title, &foundPermission.CreatedAt, &foundPermission.UpdatedAt, &foundPermission.DeletedAt)

	// If the permission is not found create a new one with the given value otherwise add the new value to the existing map
	if err != nil {
		return nil, err
	}

	// Update the permission
	err = r.PostgresDB.QueryRow(`
		UPDATE _permissions
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, title, created_at, updated_at, deleted_at
	`, payload.Name, payload.Title, currentTime, foundPermission.ID).Scan(&permission.ID, &permission.Name, &permission.Title, &permission.CreatedAt, &permission.UpdatedAt, &permission.DeletedAt)

	// If the permission does not update, return an error
	if err != nil {
		return nil, err
	}

	// Return the permission
	return &permission, nil
}

func (r *PermissionRepositoryPostgresDB) GetPermissionsByIds(ids []int64) ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM _permissions
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Title, &permission.CreatedAt, &permission.UpdatedAt, &permission.DeletedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}
