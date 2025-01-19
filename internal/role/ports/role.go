package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	"github.com/hoitek/Maja-Service/internal/role/models"
)

type RoleService interface {
	Query(dataModel *models.RolesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.RolesCreateRequestBody) (*domain.Role, error)
	Delete(payload *models.RolesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error)
	GetRoleByName(roleName string) *domain.Role
	GetRoleByID(id int) *domain.Role
	GetRolesByIds(ids []int64) ([]*domain.Role, error)
	GetRolesByUserID(userID int64) ([]*domain.Role, error)
}

type RoleRepositoryPostgresDB interface {
	Query(dataModel *models.RolesQueryRequestParams) ([]*domain.Role, error)
	Count(dataModel *models.RolesQueryRequestParams) (int64, error)
	Create(payload *models.RolesCreateRequestBody) (*domain.Role, error)
	Delete(payload *models.RolesDeleteRequestBody) ([]int64, error)
	Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error)
	GetRolesByIds(ids []int64) ([]*domain.Role, error)
	GetRolesByUserID(userID int64) ([]*domain.Role, error)
}
