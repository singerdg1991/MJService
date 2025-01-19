package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/geartype/constants"
	"github.com/hoitek/Maja-Service/internal/geartype/domain"
	"github.com/hoitek/Maja-Service/internal/geartype/models"
	"github.com/hoitek/Maja-Service/internal/geartype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type GearTypeService struct {
	PostgresRepository ports.GearTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewGearTypeService(pDB ports.GearTypeRepositoryPostgresDB, m *storage.MinIO) GearTypeService {
	go minio.SetupMinIOStorage(constants.GEAR_TYPE_BUCKET_NAME, m)
	return GearTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *GearTypeService) Query(q *models.GearTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying geartypes", q)
	geartypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.GearTypesQueryRequestParams{
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
	var items []*domain.GearType
	for _, item := range geartypes {
		items = append(items, &domain.GearType{
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

func (s *GearTypeService) Create(payload *models.GearTypesCreateRequestBody) (*domain.GearType, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *GearTypeService) Delete(payload *models.GearTypesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *GearTypeService) Update(payload *models.GearTypesCreateRequestBody, name string) (*domain.GearType, error) {
	return s.PostgresRepository.Update(payload, name)
}
