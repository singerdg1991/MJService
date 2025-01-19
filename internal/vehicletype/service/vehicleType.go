package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/vehicletype/constants"
	"github.com/hoitek/Maja-Service/internal/vehicletype/domain"
	"github.com/hoitek/Maja-Service/internal/vehicletype/models"
	"github.com/hoitek/Maja-Service/internal/vehicletype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type VehicleTypeService struct {
	PostgresRepository ports.VehicleTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewVehicleTypeService(pDB ports.VehicleTypeRepositoryPostgresDB, m *storage.MinIO) VehicleTypeService {
	go minio.SetupMinIOStorage(constants.VEHICLE_TYPE_BUCKET_NAME, m)
	return VehicleTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *VehicleTypeService) Query(q *models.VehicleTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying vehicletypes", q)
	vehicletypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.VehicleTypesQueryRequestParams{
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
	var items []*domain.VehicleType
	for _, item := range vehicletypes {
		items = append(items, &domain.VehicleType{
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

func (s *VehicleTypeService) FindByName(name string) (*domain.VehicleType, error) {
	r, err := s.Query(&models.VehicleTypesQueryRequestParams{
		Filters: models.VehicleTypeFilterType{
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
		return nil, errors.New("vehicle type not found")
	}
	vehicleTypes := r.Items.([]*domain.VehicleType)
	return vehicleTypes[0], nil
}

func (s *VehicleTypeService) Create(payload *models.VehicleTypesCreateRequestBody) (*domain.VehicleType, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *VehicleTypeService) Delete(payload *models.VehicleTypesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *VehicleTypeService) Update(payload *models.VehicleTypesCreateRequestBody, name string) (*domain.VehicleType, error) {
	return s.PostgresRepository.Update(payload, name)
}
