package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	"github.com/hoitek/Maja-Service/internal/role/models"
)

type RoleServiceStub struct {
}

func NewRoleServiceStub() *RoleServiceStub {
	return &RoleServiceStub{}
}

func (s *RoleServiceStub) Query(q *models.RolesQueryRequestParams) (*restypes.QueryResponse, error) {
	return nil, nil
}

func (s *RoleServiceStub) Create(payload *models.RolesCreateRequestBody) (*domain.Role, error) {
	return nil, nil
}

func (s *RoleServiceStub) Delete(payload *models.RolesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	return nil, nil
}

func (s *RoleServiceStub) Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error) {
	return nil, nil
}

func (s *RoleServiceStub) GetRoleByName(roleName string) *domain.Role {
	return nil
}

func (s *RoleServiceStub) GetRoleByID(id int) *domain.Role {
	return nil
}

func (s *RoleServiceStub) GetRolesByIds(ids []int64) ([]*domain.Role, error) {
	return nil, nil
}

func (s *RoleServiceStub) GetRolesByUserID(userID int64) ([]*domain.Role, error) {
	return nil, nil
}
