package service

import (
	"errors"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/trash/constants"
	"github.com/hoitek/Maja-Service/internal/trash/domain"
	"github.com/hoitek/Maja-Service/internal/trash/models"
	"github.com/hoitek/Maja-Service/internal/trash/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"
)

type TrashService struct {
	PostgresRepository ports.TrashRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewTrashService(pDB ports.TrashRepositoryPostgresDB, m *storage.MinIO) TrashService {
	go minio.SetupMinIOStorage(constants.TRASH_BUCKET_NAME, m)
	return TrashService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *TrashService) Query(q *models.TrashesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying trashes", q)
	trashes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.TrashesQueryRequestParams{
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
	var items []*domain.Trash
	for _, item := range trashes {
		items = append(items, &domain.Trash{
			ID:        item.ID,
			ModelName: item.ModelName,
			ModelID:   item.ModelID,
			Reason:    item.Reason,
			CreatedAt: item.CreatedAt,
			CreatedBy: item.CreatedBy,
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

func (s *TrashService) Create(payload *models.TrashesCreateRequestBody) (*domain.Trash, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *TrashService) Delete(payload *models.TrashesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *TrashService) FindByID(id int64) (*domain.Trash, error) {
	r, err := s.Query(&models.TrashesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("trash not found")
	}
	Trashes := r.Items.([]*domain.Trash)
	return Trashes[0], nil
}
