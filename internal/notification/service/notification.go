package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/notification/constants"
	"github.com/hoitek/Maja-Service/internal/notification/domain"
	"github.com/hoitek/Maja-Service/internal/notification/models"
	"github.com/hoitek/Maja-Service/internal/notification/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type NotificationService struct {
	PostgresRepository ports.NotificationRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewNotificationService(pDB ports.NotificationRepositoryPostgresDB, m *storage.MinIO) NotificationService {
	go minio.SetupMinIOStorage(constants.NOTIFICATION_BUCKET_NAME, m)
	return NotificationService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *NotificationService) Query(q *models.NotificationsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying notifications", q)
	notifications, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.NotificationsQueryRequestParams{
		ID:      q.ID,
		UserID:  q.UserID,
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
		Items:      notifications,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *NotificationService) Delete(payload *models.NotificationsDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// NOTIFICATION this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *NotificationService) GetNotificationsByIds(ids []int64) ([]*domain.Notification, error) {
	return s.PostgresRepository.GetNotificationsByIds(ids)
}

func (s *NotificationService) FindByID(id int64) (*domain.Notification, error) {
	r, err := s.Query(&models.NotificationsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("notification not found")
	}
	notifications, ok := r.Items.([]*domain.Notification)
	if !ok {
		return nil, errors.New("notification not found")
	}
	if len(notifications) < 1 {
		return nil, errors.New("notification not found")
	}
	return notifications[0], nil
}

func (s *NotificationService) FindByTitle(title string) (*domain.Notification, error) {
	r, err := s.Query(&models.NotificationsQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.NotificationFilterType{
			Title: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: title,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("notification not found")
	}
	notifications, ok := r.Items.([]*domain.Notification)
	if !ok {
		return nil, errors.New("notification not found")
	}
	if len(notifications) < 1 {
		return nil, errors.New("notification not found")
	}
	return notifications[0], nil
}

func (s *NotificationService) CreateNotification(userId int64, title, description string, status string, extra interface{}) (*domain.Notification, error) {
	return s.PostgresRepository.CreateNotification(userId, title, description, status, false, extra)
}

func (s *NotificationService) CreateSystemNotification(userId int64, title, description string, status string, extra interface{}) (*domain.Notification, error) {
	return s.PostgresRepository.CreateNotification(userId, title, description, status, true, extra)
}

func (s *NotificationService) UpdateStatus(payload *models.NotificationsUpdateStatusRequestBody, notificationID int64) (*domain.Notification, error) {
	return s.PostgresRepository.UpdateStatus(payload, notificationID)
}
