package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	"github.com/hoitek/Maja-Service/internal/role/models"
	"github.com/hoitek/Maja-Service/utils"
	"github.com/lib/pq"
)

type RoleRepositoryPostgresDB struct {
	PostgresDB *sql.DB
}

func NewRoleRepositoryPostgresDB(d *sql.DB) *RoleRepositoryPostgresDB {
	return &RoleRepositoryPostgresDB{
		PostgresDB: d,
	}
}

func makeWhereFilters(queries *models.RolesQueryRequestParams) []string {
	var where []string
	if queries != nil {
		if queries.ID != 0 {
			where = append(where, fmt.Sprintf(" r.id = %d", queries.ID))
		}
		if queries.Filters.Name.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Name.Op, fmt.Sprintf("%v", queries.Filters.Name.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.name %s %s", opValue.Operator, val))
		}
		if queries.Filters.Type.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.Type.Op, fmt.Sprintf("%v", queries.Filters.Type.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.type %s %s", opValue.Operator, val))
		}
		if queries.Filters.CreatedAt.Op != "" {
			opValue := utils.GetDBOperatorAndValue(queries.Filters.CreatedAt.Op, fmt.Sprintf("%v", queries.Filters.CreatedAt.Value))
			val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
			where = append(where, fmt.Sprintf(" r.created_at %s %s", opValue.Operator, val))
		}
	}
	return where
}

func (r *RoleRepositoryPostgresDB) Query(queries *models.RolesQueryRequestParams) ([]*domain.Role, error) {
	q := `SELECT * FROM _roles r `
	if queries != nil {
		where := makeWhereFilters(queries)
		if len(where) > 0 {
			q += " WHERE " + strings.Join(where, " AND ")
		}
		var sorts []string
		if queries.Sorts.Name.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" r.name %s", queries.Sorts.Name.Op))
		}
		if queries.Sorts.Type.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" r.type %s", queries.Sorts.Type.Op))
		}
		if queries.Sorts.CreatedAt.Op != "" {
			sorts = append(sorts, fmt.Sprintf(" r.created_at %s", queries.Sorts.CreatedAt.Op))
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

	var roles []*domain.Role
	rows, err := r.PostgresDB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Type,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Find Permissions
		pRows, err := r.PostgresDB.Query(`
			SELECT p.id, p.name, p.title
			FROM _permissions p
			INNER JOIN _rolesPermissions rp ON rp.permissionId = p.id
			WHERE rp.roleId = $1
		`, role.ID)
		if err != nil {
			pRows.Close()
			return nil, err
		}
		var permissions []*domain.RolePermission
		for pRows.Next() {
			var permission domain.RolePermission
			err := pRows.Scan(
				&permission.ID,
				&permission.Name,
				&permission.Title,
			)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			permissions = append(permissions, &permission)
		}
		pRows.Close()
		role.Permissions = permissions
		roles = append(roles, &role)
	}
	return roles, nil
}

func (r *RoleRepositoryPostgresDB) Count(queries *models.RolesQueryRequestParams) (int64, error) {
	q := `SELECT COUNT(*) FROM _roles r `
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

func (r *RoleRepositoryPostgresDB) Create(payload *models.RolesCreateRequestBody) (*domain.Role, error) {
	var role domain.Role

	// Current time
	currentTime := time.Now()

	// Find the role by name
	var foundRole domain.Role
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM _roles
		WHERE name = $1
	`, payload.Name).Scan(&foundRole.ID, &foundRole.Name, &foundRole.Type, &foundRole.CreatedAt, &foundRole.UpdatedAt, &foundRole.DeletedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Insert the role
	var roleType = constants.ROLE_TYPE_INTERNAL
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(`
		INSERT INTO _roles (name, type, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, type, created_at, updated_at, deleted_at
	`, payload.Name, roleType, currentTime, currentTime, nil).Scan(&role.ID, &role.Name, &role.Type, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, permissionID := range payload.PermissionsInt64 {
		_, err = tx.Exec(`
			INSERT INTO _rolesPermissions (roleId, permissionId, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5)
		`, role.ID, permissionID, currentTime, currentTime, nil)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	roles, err := r.Query(&models.RolesQueryRequestParams{ID: int(role.ID)})
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.New("role not found")
	}

	// Return the role
	return roles[0], nil
}

func (r *RoleRepositoryPostgresDB) Delete(payload *models.RolesDeleteRequestBody) ([]int64, error) {
	var rowsAffected int64
	err := r.PostgresDB.QueryRow(`
		DELETE FROM _roles
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

func (r *RoleRepositoryPostgresDB) Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error) {
	var role domain.Role

	// Current time
	currentTime := time.Now()

	// Find the role by id
	var foundRole domain.Role
	err := r.PostgresDB.QueryRow(`
		SELECT *
		FROM _roles
		WHERE id = $1
	`, id).Scan(&foundRole.ID, &foundRole.Name, &foundRole.Type, &foundRole.CreatedAt, &foundRole.UpdatedAt, &foundRole.DeletedAt)
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx, err := r.PostgresDB.Begin()
	if err != nil {
		return nil, err
	}

	// Update the role
	err = tx.QueryRow(`
		UPDATE _roles
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, type, created_at, updated_at, deleted_at
	`, payload.Name, currentTime, foundRole.ID).Scan(&role.ID, &role.Name, &role.Type, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete all permissions
	_, err = tx.Exec(`
		DELETE FROM _rolesPermissions
		WHERE roleId = $1
	`, foundRole.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert new permissions
	for _, permissionID := range payload.PermissionsInt64 {
		_, err = tx.Exec(`
			INSERT INTO _rolesPermissions (roleId, permissionId, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5)
		`, foundRole.ID, permissionID, currentTime, currentTime, nil)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Query the role
	roles, err := r.Query(&models.RolesQueryRequestParams{ID: int(role.ID)})
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.New("role not found")
	}

	// Return the role
	return roles[0], nil
}

func (r *RoleRepositoryPostgresDB) GetRolesByIds(ids []int64) ([]*domain.Role, error) {
	var roles []*domain.Role
	rows, err := r.PostgresDB.Query(`
		SELECT *
		FROM _roles
		WHERE id = ANY ($1)
	`, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Type,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Find Permissions
		pRows, err := r.PostgresDB.Query(`
			SELECT p.id, p.name, p.title
			FROM _permissions p
			INNER JOIN _rolesPermissions rp ON rp.permissionId = p.id
			WHERE rp.roleId = $1
		`, role.ID)
		if err != nil {
			pRows.Close()
			return nil, err
		}
		var permissions []*domain.RolePermission
		for pRows.Next() {
			var permission domain.RolePermission
			err := pRows.Scan(
				&permission.ID,
				&permission.Name,
				&permission.Title,
			)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			permissions = append(permissions, &permission)
		}
		pRows.Close()
		role.Permissions = permissions
		roles = append(roles, &role)
	}
	return roles, nil
}

func (r *RoleRepositoryPostgresDB) GetRolesByUserID(userID int64) ([]*domain.Role, error) {
	var roles []*domain.Role
	rows, err := r.PostgresDB.Query(`
        SELECT r.id, r.name, r.type, r.created_at, r.updated_at, r.deleted_at
        FROM _roles r
        INNER JOIN usersRoles ur ON ur.roleId = r.id
        WHERE ur.userId = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Type,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Find Permissions
		pRows, err := r.PostgresDB.Query(`
            SELECT p.id, p.name, p.title
            FROM _permissions p
            INNER JOIN _rolesPermissions rp ON rp.permissionId = p.id
            WHERE rp.roleId = $1
        `, role.ID)
		if err != nil {
			pRows.Close()
			return nil, err
		}
		var permissions []*domain.RolePermission
		for pRows.Next() {
			var permission domain.RolePermission
			err := pRows.Scan(
				&permission.ID,
				&permission.Name,
				&permission.Title,
			)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			permissions = append(permissions, &permission)
		}
		role.Permissions = permissions
		roles = append(roles, &role)
		pRows.Close()
	}
	return roles, nil
}
