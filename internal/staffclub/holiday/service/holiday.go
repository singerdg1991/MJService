package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/constants"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type HolidayService struct {
	PostgresRepository ports.HolidayRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewHolidayService(pDB ports.HolidayRepositoryPostgresDB, m *storage.MinIO) HolidayService {
	go minio.SetupMinIOStorage(constants.HOLIDAY_BUCKET_NAME, m)
	return HolidayService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *HolidayService) Query(q *models.HolidaysQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying holidays", q)
	holidays, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.HolidaysQueryRequestParams{
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

	return &restypes.QueryResponse{
		Items:      holidays,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *HolidayService) Create(payload *models.HolidaysCreateRequestBody) (*domain.Holiday, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *HolidayService) Delete(payload *models.HolidaysDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *HolidayService) Update(payload *models.HolidaysCreateRequestBody, id int64) (*domain.Holiday, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *HolidayService) GetHolidaysByIds(ids []int64) ([]*domain.Holiday, error) {
	return s.PostgresRepository.GetHolidaysByIds(ids)
}

func (s *HolidayService) FindByID(id int64) (*domain.Holiday, error) {
	r, err := s.Query(&models.HolidaysQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("holiday not found")
	}
	holidays := r.Items.([]*domain.Holiday)
	return holidays[0], nil
}

func (s *HolidayService) UpdateStatus(payload *models.HolidaysUpdateStatusRequestBody, id int64) (*domain.Holiday, error) {
	return s.PostgresRepository.UpdateStatus(payload, id)
}
