package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/permission/constants"
	"github.com/hoitek/Maja-Service/internal/permission/domain"
	"github.com/hoitek/Maja-Service/internal/permission/models"
	"github.com/hoitek/Maja-Service/internal/permission/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type PermissionService struct {
	PostgresRepository ports.PermissionRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewPermissionService(pDB ports.PermissionRepositoryPostgresDB, m *storage.MinIO) PermissionService {
	go minio.SetupMinIOStorage(constants.PERMISSION_BUCKET_NAME, m)
	return PermissionService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *PermissionService) Query(q *models.PermissionsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying permissions", q)
	permissions, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.PermissionsQueryRequestParams{
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

	// Transform the response to the format that the frontend expects
	var items []*domain.Permission
	for _, item := range permissions {
		items = append(items, &domain.Permission{
			ID:    item.ID,
			Name:  item.Name,
			Title: item.Title,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *PermissionService) Create(payload *models.PermissionsCreateRequestBody) (*domain.Permission, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *PermissionService) Delete(payload *models.PermissionsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *PermissionService) Update(payload *models.PermissionsCreateRequestBody, id int64) (*domain.Permission, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *PermissionService) GetPermissionsByIds(ids []int64) ([]*domain.Permission, error) {
	return s.PostgresRepository.GetPermissionsByIds(ids)
}

func (s *PermissionService) FindByID(id int64) (*domain.Permission, error) {
	r, err := s.Query(&models.PermissionsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	permissions := r.Items.([]*domain.Permission)
	return permissions[0], nil
}

func (s *PermissionService) FindByName(name string) (*domain.Permission, error) {
	r, err := s.Query(&models.PermissionsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.PermissionFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	permissions := r.Items.([]*domain.Permission)
	return permissions[0], nil
}
