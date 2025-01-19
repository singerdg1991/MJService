package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/notification/domain"
	"github.com/hoitek/Maja-Service/internal/notification/models"
)

type NotificationRepositoryStub struct {
	Notifications []*domain.Notification
}

type notificationTestCondition struct {
	HasError bool
}

var UserTestCondition *notificationTestCondition = &notificationTestCondition{}

func NewNotificationRepositoryStub() *NotificationRepositoryStub {
	return &NotificationRepositoryStub{
		Notifications: []*domain.Notification{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *NotificationRepositoryStub) Query(dataModel *models.NotificationsQueryRequestParams) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	for _, v := range r.Notifications {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			notifications = append(notifications, v)
			break
		}
	}
	return notifications, nil
}

func (r *NotificationRepositoryStub) Count(dataModel *models.NotificationsQueryRequestParams) (int64, error) {
	var notifications []*domain.Notification
	for _, v := range r.Notifications {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			notifications = append(notifications, v)
			break
		}
	}
	return int64(len(notifications)), nil
}

func (r *NotificationRepositoryStub) Delete(payload *models.NotificationsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *NotificationRepositoryStub) GetNotificationsByIds(ids []int64) ([]*domain.Notification, error) {
	panic("implement me")
}

func (r *NotificationRepositoryStub) CreateNotification(userId int64, title, description string, status string, isForSystem bool, extra interface{}) (*domain.Notification, error) {
	panic("implement me")
}

func (r *NotificationRepositoryStub) UpdateStatus(payload *models.NotificationsUpdateStatusRequestBody, notificationID int64) (*domain.Notification, error) {
	panic("implement me")
}
