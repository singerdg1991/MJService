package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/punishment/constants"
	"github.com/hoitek/Maja-Service/internal/punishment/domain"
	"github.com/hoitek/Maja-Service/internal/punishment/models"
	"github.com/hoitek/Maja-Service/internal/punishment/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type PunishmentService struct {
	PostgresRepository ports.PunishmentRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewPunishmentService(pDB ports.PunishmentRepositoryPostgresDB, m *storage.MinIO) PunishmentService {
	go minio.SetupMinIOStorage(constants.PUNISHMENT_BUCKET_NAME, m)
	return PunishmentService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *PunishmentService) Query(q *models.PunishmentsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying punishments", q)
	punishments, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.PunishmentsQueryRequestParams{
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
	var items []*domain.Punishment
	for _, item := range punishments {
		items = append(items, &domain.Punishment{
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

func (s *PunishmentService) Create(payload *models.PunishmentsCreateRequestBody) (*domain.Punishment, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *PunishmentService) Delete(payload *models.PunishmentsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *PunishmentService) Update(payload *models.PunishmentsCreateRequestBody, id int64) (*domain.Punishment, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *PunishmentService) GetPunishmentsByIds(ids []int64) ([]*domain.Punishment, error) {
	return s.PostgresRepository.GetPunishmentsByIds(ids)
}

func (s *PunishmentService) FindByID(id int64) (*domain.Punishment, error) {
	r, err := s.Query(&models.PunishmentsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	punishments := r.Items.([]*domain.Punishment)
	return punishments[0], nil
}

func (s *PunishmentService) FindByName(name string) (*domain.Punishment, error) {
	r, err := s.Query(&models.PunishmentsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.PunishmentFilterType{
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
	punishments := r.Items.([]*domain.Punishment)
	return punishments[0], nil
}
