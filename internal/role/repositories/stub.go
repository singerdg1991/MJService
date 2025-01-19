package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	"github.com/hoitek/Maja-Service/internal/role/models"
)

type RoleRepositoryStub struct {
	Roles []*domain.Role
}

type roleTestCondition struct {
	HasError bool
}

var UserTestCondition *roleTestCondition = &roleTestCondition{}

func NewRoleRepositoryStub() *RoleRepositoryStub {
	return &RoleRepositoryStub{
		Roles: []*domain.Role{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *RoleRepositoryStub) Query(dataModel *models.RolesQueryRequestParams) ([]*domain.Role, error) {
	var roles []*domain.Role
	for _, v := range r.Roles {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			roles = append(roles, v)
			break
		}
	}
	return roles, nil
}

func (r *RoleRepositoryStub) Count(dataModel *models.RolesQueryRequestParams) (int64, error) {
	var roles []*domain.Role
	for _, v := range r.Roles {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			roles = append(roles, v)
			break
		}
	}
	return int64(len(roles)), nil
}

func (r *RoleRepositoryStub) Migrate() {
	// do stuff
}

func (r *RoleRepositoryStub) Seed() {
	// do stuff
}

func (r *RoleRepositoryStub) Create(payload *models.RolesCreateRequestBody) (*domain.Role, error) {
	panic("implement me")
}

func (r *RoleRepositoryStub) Delete(payload *models.RolesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *RoleRepositoryStub) Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error) {
	panic("implement me")
}

func (r *RoleRepositoryStub) GetRolesByIds(ids []int64) ([]*domain.Role, error) {
	panic("implement me")
}

func (r *RoleRepositoryStub) GetRolesByUserID(userID int64) ([]*domain.Role, error) {
	panic("implement me")
}
