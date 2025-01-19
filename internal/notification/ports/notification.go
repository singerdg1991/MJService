package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/notification/domain"
	"github.com/hoitek/Maja-Service/internal/notification/models"
)

type NotificationService interface {
	Query(dataModel *models.NotificationsQueryRequestParams) (*restypes.QueryResponse, error)
	Delete(payload *models.NotificationsDeleteRequestBody) (*restypes.DeleteResponse, error)
	GetNotificationsByIds(ids []int64) ([]*domain.Notification, error)
	FindByID(id int64) (*domain.Notification, error)
	FindByTitle(title string) (*domain.Notification, error)
	CreateNotification(userId int64, title, description string, status string, extra interface{}) (*domain.Notification, error)
	CreateSystemNotification(userId int64, title, description string, status string, extra interface{}) (*domain.Notification, error)
	UpdateStatus(payload *models.NotificationsUpdateStatusRequestBody, notificationID int64) (*domain.Notification, error)
}

type NotificationRepositoryPostgresDB interface {
	Query(dataModel *models.NotificationsQueryRequestParams) ([]*domain.Notification, error)
	Count(dataModel *models.NotificationsQueryRequestParams) (int64, error)
	Delete(payload *models.NotificationsDeleteRequestBody) ([]int64, error)
	GetNotificationsByIds(ids []int64) ([]*domain.Notification, error)
	CreateNotification(userId int64, title, description string, status string, isForSystem bool, extra interface{}) (*domain.Notification, error)
	UpdateStatus(payload *models.NotificationsUpdateStatusRequestBody, notificationID int64) (*domain.Notification, error)
}
