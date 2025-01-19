package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/city/constants"
	"github.com/hoitek/Maja-Service/internal/city/domain"
	"github.com/hoitek/Maja-Service/internal/city/models"
	"github.com/hoitek/Maja-Service/internal/city/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type CityService struct {
	PostgresRepository ports.CityRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewCityService(pDB ports.CityRepositoryPostgresDB, m *storage.MinIO) CityService {
	go minio.SetupMinIOStorage(constants.CITY_BUCKET_NAME, m)
	return CityService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *CityService) Query(q *models.CitiesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying cities", q)

	cities, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.CitiesQueryRequestParams{
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
	var items []*domain.City
	for _, item := range cities {
		items = append(items, &domain.City{
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

func (s *CityService) Create(payload *models.CitiesCreateRequestBody) (*domain.City, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *CityService) Delete(payload *models.CitiesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *CityService) Update(payload *models.CitiesCreateRequestBody, name string) (*domain.City, error) {
	return s.PostgresRepository.Update(payload, name)
}

func (s *CityService) GetCityByID(id int64) (*domain.City, error) {
	cities, err := s.PostgresRepository.Query(&models.CitiesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, nil
	}
	return cities[0], nil
}
