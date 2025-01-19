package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/push/constants"
	"github.com/hoitek/Maja-Service/internal/push/domain"
	"github.com/hoitek/Maja-Service/internal/push/models"
	"github.com/hoitek/Maja-Service/internal/push/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type PushService struct {
	PostgresRepository ports.PushRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewPushService(pDB ports.PushRepositoryPostgresDB, m *storage.MinIO) PushService {
	go minio.SetupMinIOStorage(constants.PUSH_BUCKET_NAME, m)
	return PushService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *PushService) Query(q *models.PushesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying pushes", q)
	pushes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.PushesQueryRequestParams{
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
	var items []*domain.Push
	for _, item := range pushes {
		items = append(items, &domain.Push{
			ID:         item.ID,
			UserID:     item.UserID,
			Endpoint:   item.Endpoint,
			KeysAuth:   item.KeysAuth,
			KeysP256dh: item.KeysP256dh,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			DeletedAt:  item.DeletedAt,
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

func (s *PushService) Create(payload *models.PushesCreateRequestBody) (*domain.Push, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *PushService) GetPushesByIds(ids []int64) ([]*domain.Push, error) {
	return s.PostgresRepository.GetPushesByIds(ids)
}

func (s *PushService) FindByID(id int64) (*domain.Push, error) {
	r, err := s.Query(&models.PushesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("push not found")
	}
	pushes, ok := r.Items.([]*domain.Push)
	if !ok {
		return nil, errors.New("push not found")
	}
	if len(pushes) < 1 {
		return nil, errors.New("push not found")
	}
	return pushes[0], nil
}

func (s *PushService) FindByUserID(userID int) (*domain.Push, error) {
	r, err := s.Query(&models.PushesQueryRequestParams{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("push not found")
	}
	pushes, ok := r.Items.([]*domain.Push)
	if !ok {
		return nil, errors.New("push not found")
	}
	if len(pushes) < 1 {
		return nil, errors.New("push not found")
	}
	return pushes[0], nil
}
