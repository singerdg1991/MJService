package authz

import (
	"github.com/hoitek/Maja-Service/internal/_shared/authz/ports"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
)

// AuthZ is the interface that wraps the methods required for authorization
type AuthZ struct {
	RoleService rPorts.RoleService
}

// NewAuthZ returns a new AuthZ
func NewAuthZ(roleService rPorts.RoleService) ports.AuthZ {
	return &AuthZ{
		RoleService: roleService,
	}
}

// GetRolesByUserID returns a list of roles by user id
func (a *AuthZ) GetRolesByUserID(userID int64) ([]*domain.Role, error) {
	return a.RoleService.GetRolesByUserID(userID)
}

// UserHasRole returns true if the user has the given role
func (a *AuthZ) UserHasRole(userID int64, roleName string) (bool, error) {
	roles, err := a.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == roleName {
			return true, nil
		}
	}
	return false, nil
}

// UserHasPermission returns true if the user has the given permission
func (a *AuthZ) UserHasPermission(userID int64, permissionName string) (bool, error) {
	roles, err := a.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true, nil
			}
		}
	}
	return false, nil
}

// UserHasAnyPermission returns true if the user has any of the given permissions
func (a *AuthZ) UserHasAnyPermission(userID int64, permissionNames []string) (bool, error) {
	roles, err := a.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		for _, permission := range role.Permissions {
			for _, permissionName := range permissionNames {
				if permission.Name == permissionName {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// UserHasAnyRole returns true if the user has any of the given roles
func (a *AuthZ) UserHasAnyRole(userID int64, roleNames []string) (bool, error) {
	roles, err := a.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		for _, roleName := range roleNames {
			if role.Name == roleName {
				return true, nil
			}
		}
	}
	return false, nil
}
