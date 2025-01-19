package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/permission/domain"
	"github.com/hoitek/Maja-Service/internal/permission/models"
)

type PermissionService interface {
	Query(dataModel *models.PermissionsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.PermissionsCreateRequestBody) (*domain.Permission, error)
	Delete(payload *models.PermissionsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.PermissionsCreateRequestBody, id int64) (*domain.Permission, error)
	GetPermissionsByIds(ids []int64) ([]*domain.Permission, error)
	FindByID(id int64) (*domain.Permission, error)
	FindByName(name string) (*domain.Permission, error)
}

type PermissionRepositoryPostgresDB interface {
	Query(dataModel *models.PermissionsQueryRequestParams) ([]*domain.Permission, error)
	Count(dataModel *models.PermissionsQueryRequestParams) (int64, error)
	Create(payload *models.PermissionsCreateRequestBody) (*domain.Permission, error)
	Delete(payload *models.PermissionsDeleteRequestBody) ([]int64, error)
	Update(payload *models.PermissionsCreateRequestBody, id int64) (*domain.Permission, error)
	GetPermissionsByIds(ids []int64) ([]*domain.Permission, error)
}
