package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/license/constants"
	"github.com/hoitek/Maja-Service/internal/license/domain"
	"github.com/hoitek/Maja-Service/internal/license/models"
	"github.com/hoitek/Maja-Service/internal/license/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type LicenseService struct {
	PostgresRepository ports.LicenseRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewLicenseService(pDB ports.LicenseRepositoryPostgresDB, m *storage.MinIO) LicenseService {
	go minio.SetupMinIOStorage(constants.LICENSE_BUCKET_NAME, m)
	return LicenseService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *LicenseService) Query(q *models.LicensesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying licenses", q)
	licenses, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.LicensesQueryRequestParams{
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
	var items []*domain.License
	for _, item := range licenses {
		items = append(items, &domain.License{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
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

func (s *LicenseService) Create(payload *models.LicensesCreateRequestBody) (*domain.License, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *LicenseService) Delete(payload *models.LicensesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *LicenseService) Update(payload *models.LicensesCreateRequestBody, id int64) (*domain.License, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *LicenseService) GetLicensesByIds(ids []int64) ([]*domain.License, error) {
	return s.PostgresRepository.GetLicensesByIds(ids)
}

func (s *LicenseService) FindByID(id int64) (*domain.License, error) {
	r, err := s.Query(&models.LicensesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	licenses := r.Items.([]*domain.License)
	return licenses[0], nil
}

func (s *LicenseService) FindByName(name string) (*domain.License, error) {
	r, err := s.Query(&models.LicensesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.LicenseFilterType{
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
	licenses := r.Items.([]*domain.License)
	return licenses[0], nil
}
