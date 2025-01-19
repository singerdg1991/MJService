package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/shifttype/constants"
	"github.com/hoitek/Maja-Service/internal/shifttype/domain"
	"github.com/hoitek/Maja-Service/internal/shifttype/models"
	"github.com/hoitek/Maja-Service/internal/shifttype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type ShiftTypeService struct {
	PostgresRepository ports.ShiftTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewShiftTypeService(pDB ports.ShiftTypeRepositoryPostgresDB, m *storage.MinIO) ShiftTypeService {
	go minio.SetupMinIOStorage(constants.SHIFT_TYPE_BUCKET_NAME, m)
	return ShiftTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ShiftTypeService) Query(q *models.ShiftTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying ShiftTypes", q)
	shiftTypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ShiftTypesQueryRequestParams{
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
	var items []*domain.ShiftType
	for _, item := range shiftTypes {
		items = append(items, &domain.ShiftType{
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

func (s *ShiftTypeService) GetShiftTypesByIds(ids []int64) ([]*domain.ShiftType, error) {
	return s.PostgresRepository.GetShiftTypesByIds(ids)
}
