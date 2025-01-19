package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/equipment/constants"
	"github.com/hoitek/Maja-Service/internal/equipment/domain"
	"github.com/hoitek/Maja-Service/internal/equipment/models"
	"github.com/hoitek/Maja-Service/internal/equipment/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type EquipmentService struct {
	PostgresRepository ports.EquipmentRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewEquipmentService(pDB ports.EquipmentRepositoryPostgresDB, m *storage.MinIO) EquipmentService {
	go minio.SetupMinIOStorage(constants.EQUIPMENT_BUCKET_NAME, m)
	return EquipmentService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *EquipmentService) Query(q *models.EquipmentsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying equipments", q)
	equipments, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.EquipmentsQueryRequestParams{
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
	var items []*domain.Equipment
	for _, item := range equipments {
		items = append(items, &domain.Equipment{
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

func (s *EquipmentService) Create(payload *models.EquipmentsCreateRequestBody) (*domain.Equipment, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *EquipmentService) Delete(payload *models.EquipmentsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *EquipmentService) Update(payload *models.EquipmentsCreateRequestBody, id int64) (*domain.Equipment, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *EquipmentService) GetEquipmentsByIds(ids []int64) ([]*domain.Equipment, error) {
	return s.PostgresRepository.GetEquipmentsByIds(ids)
}

func (s *EquipmentService) FindByID(id int64) (*domain.Equipment, error) {
	r, err := s.Query(&models.EquipmentsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	equipments := r.Items.([]*domain.Equipment)
	return equipments[0], nil
}

func (s *EquipmentService) FindByName(name string) (*domain.Equipment, error) {
	r, err := s.Query(&models.EquipmentsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.EquipmentFilterType{
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
	equipments := r.Items.([]*domain.Equipment)
	return equipments[0], nil
}
