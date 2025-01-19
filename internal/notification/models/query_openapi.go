package models

import (
	"github.com/hoitek/Maja-Service/internal/notification/domain"
)

/*
 * @apiDefine: NotificationsResponseData
 */
type NotificationsResponseData struct {
	ID          uint                    `json:"id" openapi:"example:1"`
	UserID      uint                    `json:"userId" openapi:"example:1"`
	User        domain.NotificationUser `json:"user" openapi:"$ref:NotificationUser;type:object;"`
	Title       string                  `json:"title" openapi:"example:title;required"`
	Description string                  `json:"description" openapi:"example:description;required"`
	ReadAt      *string                 `json:"read_at" openapi:"example:2021-01-01T00:00:00Z"`
	ReadBy      domain.NotificationUser `json:"readBy" openapi:"$ref:NotificationUser;type:object;"`
	IsForSystem bool                    `json:"isForSystem" openapi:"example:false"`
	Status      *string                 `json:"status" openapi:"example:pending"`
	StatusAt    *string                 `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt   string                  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: NotificationsQueryResponseData
 */
type NotificationsQueryResponseData struct {
	Limit      int                         `json:"limit" openapi:"example:10"`
	Offset     int                         `json:"offset" openapi:"example:0"`
	Page       int                         `json:"page" openapi:"example:1"`
	TotalRows  int                         `json:"totalRows" openapi:"example:1"`
	TotalPages int                         `json:"totalPages" openapi:"example:1"`
	Items      []NotificationsResponseData `json:"items" openapi:"$ref:NotificationsResponseData;type:array"`
}

/*
 * @apiDefine: NotificationsQueryResponse
 */
type NotificationsQueryResponse struct {
	StatusCode int                            `json:"statusCode" openapi:"example:200;"`
	Data       NotificationsQueryResponseData `json:"data" openapi:"$ref:NotificationsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: NotificationsQueryNotFoundResponse
 */
type NotificationsQueryNotFoundResponse struct {
	Notifications []domain.Notification `json:"notifications" openapi:"$ref:Notification;type:array"`
}
