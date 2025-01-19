package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/stafftype/constants"
	"github.com/hoitek/Maja-Service/internal/stafftype/domain"
	"github.com/hoitek/Maja-Service/internal/stafftype/models"
	"github.com/hoitek/Maja-Service/internal/stafftype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type StaffTypeService struct {
	PostgresRepository ports.StaffTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewStaffTypeService(pDB ports.StaffTypeRepositoryPostgresDB, m *storage.MinIO) StaffTypeService {
	go minio.SetupMinIOStorage(constants.STAFF_TYPE_BUCKET_NAME, m)
	return StaffTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *StaffTypeService) Query(q *models.StaffTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying staffTypes", q)
	staffTypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.StaffTypesQueryRequestParams{
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
	var items []*domain.StaffType
	for _, item := range staffTypes {
		items = append(items, &domain.StaffType{
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

func (s *StaffTypeService) GetStaffTypesByIds(ids []int64) ([]*domain.StaffType, error) {
	return s.PostgresRepository.GetStaffTypesByIds(ids)
}
