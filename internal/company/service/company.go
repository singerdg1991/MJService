package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/company/constants"
	"github.com/hoitek/Maja-Service/internal/company/domain"
	"github.com/hoitek/Maja-Service/internal/company/models"
	"github.com/hoitek/Maja-Service/internal/company/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type CompanyService struct {
	PostgresRepository ports.CompanyRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewCompanyService(pDB ports.CompanyRepositoryPostgresDB, m *storage.MinIO) CompanyService {
	go minio.SetupMinIOStorage(constants.COMPANY_BUCKET_NAME, m)
	return CompanyService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *CompanyService) Query(q *models.CompaniesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying companies", q)
	companies, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.CompaniesQueryRequestParams{
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
	var items []*domain.Company
	for _, item := range companies {
		items = append(items, &domain.Company{
			ID:   item.ID,
			Name: item.Name,
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

func (s *CompanyService) Create(payload *models.CompaniesCreateRequestBody) (*domain.Company, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *CompanyService) Delete(payload *models.CompaniesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *CompanyService) Update(payload *models.CompaniesCreateRequestBody, name string) (*domain.Company, error) {
	return s.PostgresRepository.Update(payload, name)
}
