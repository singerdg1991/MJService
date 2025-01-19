package service

import (
	"database/sql"
	"errors"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/vehicle/constants"
	"github.com/hoitek/Maja-Service/internal/vehicle/domain"
	"github.com/hoitek/Maja-Service/internal/vehicle/models"
	"github.com/hoitek/Maja-Service/internal/vehicle/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type VehicleService struct {
	PostgresRepository ports.VehicleRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewVehicleService(pDB ports.VehicleRepositoryPostgresDB, m *storage.MinIO) VehicleService {
	go minio.SetupMinIOStorage(constants.VEHICLE_BUCKET_NAME, m)
	return VehicleService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *VehicleService) Query(q *models.VehiclesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying vehicles", q)
	vehicles, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.VehiclesQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			count = 0
		} else {
			return nil, err
		}
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
		Items:      vehicles,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *VehicleService) Create(payload *models.VehiclesCreateRequestBody) (*domain.Vehicle, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *VehicleService) Delete(payload *models.VehiclesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *VehicleService) Update(payload *models.VehiclesCreateRequestBody, id int) (*domain.Vehicle, error) {
	return s.PostgresRepository.Update(payload, id)
}
