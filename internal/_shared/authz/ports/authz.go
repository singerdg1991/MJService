package ports

import "github.com/hoitek/Maja-Service/internal/role/domain"

type AuthZ interface {
	GetRolesByUserID(userID int64) ([]*domain.Role, error)
	UserHasRole(userID int64, roleName string) (bool, error)
	UserHasPermission(userID int64, permissionName string) (bool, error)
	UserHasAnyPermission(userID int64, permissionNames []string) (bool, error)
	UserHasAnyRole(userID int64, roleNames []string) (bool, error)
}
