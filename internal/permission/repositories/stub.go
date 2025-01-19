package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/permission/domain"
	"github.com/hoitek/Maja-Service/internal/permission/models"
)

type PermissionRepositoryStub struct {
	Permissions []*domain.Permission
}

type permissionTestCondition struct {
	HasError bool
}

var UserTestCondition *permissionTestCondition = &permissionTestCondition{}

func NewPermissionRepositoryStub() *PermissionRepositoryStub {
	return &PermissionRepositoryStub{
		Permissions: []*domain.Permission{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *PermissionRepositoryStub) Query(dataModel *models.PermissionsQueryRequestParams) ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	for _, v := range r.Permissions {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			permissions = append(permissions, v)
			break
		}
	}
	return permissions, nil
}

func (r *PermissionRepositoryStub) Count(dataModel *models.PermissionsQueryRequestParams) (int64, error) {
	var permissions []*domain.Permission
	for _, v := range r.Permissions {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			permissions = append(permissions, v)
			break
		}
	}
	return int64(len(permissions)), nil
}

func (r *PermissionRepositoryStub) Migrate() {
	// do stuff
}

func (r *PermissionRepositoryStub) Seed() {
	// do stuff
}

func (r *PermissionRepositoryStub) Create(payload *models.PermissionsCreateRequestBody) (*domain.Permission, error) {
	panic("implement me")
}

func (r *PermissionRepositoryStub) Delete(payload *models.PermissionsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *PermissionRepositoryStub) Update(payload *models.PermissionsCreateRequestBody, id int64) (*domain.Permission, error) {
	panic("implement me")
}

func (r *PermissionRepositoryStub) GetPermissionsByIds(ids []int64) ([]*domain.Permission, error) {
	panic("implement me")
}
