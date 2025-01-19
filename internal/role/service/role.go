package service

import (
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/role/constants"
	"github.com/hoitek/Maja-Service/internal/role/domain"
	"github.com/hoitek/Maja-Service/internal/role/models"
	"github.com/hoitek/Maja-Service/internal/role/ports"
	"github.com/hoitek/Maja-Service/storage"
	"math"

	"github.com/hoitek/Kit/exp"
)

type RoleService struct {
	PostgresRepository ports.RoleRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewRoleService(pDB ports.RoleRepositoryPostgresDB, m *storage.MinIO) *RoleService {
	go minio.SetupMinIOStorage(constants.ROLE_BUCKET_NAME, m)
	return &RoleService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *RoleService) Query(q *models.RolesQueryRequestParams) (*restypes.QueryResponse, error) {
	roles, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.RolesQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      roles,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *RoleService) Create(payload *models.RolesCreateRequestBody) (*domain.Role, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *RoleService) Delete(payload *models.RolesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// TODO this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *RoleService) Update(payload *models.RolesCreateRequestBody, id int64) (*domain.Role, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *RoleService) GetRoleByName(roleName string) *domain.Role {
	r, err := s.Query(&models.RolesQueryRequestParams{
		Filters: models.RoleFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: roleName,
			},
		},
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	roles, ok := r.Items.([]*domain.Role)
	if !ok {
		return nil
	}
	return roles[0]
}

func (s *RoleService) GetRoleByID(id int) *domain.Role {
	r, err := s.Query(&models.RolesQueryRequestParams{
		ID: id,
	})
	if err != nil {
		return nil
	}
	if r.TotalRows <= 0 {
		return nil
	}
	roles, ok := r.Items.([]*domain.Role)
	if !ok {
		return nil
	}
	return roles[0]
}

func (s *RoleService) GetRolesByIds(ids []int64) ([]*domain.Role, error) {
	return s.PostgresRepository.GetRolesByIds(ids)
}

func (s *RoleService) GetRolesByUserID(userID int64) ([]*domain.Role, error) {
	return s.PostgresRepository.GetRolesByUserID(userID)
}
