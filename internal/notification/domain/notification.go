package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: NotificationUser
 */
type NotificationUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.png"`
}

/*
 * @apiDefine: Notification
 */
type Notification struct {
	ID          uint              `json:"id" openapi:"example:1"`
	UserID      *uint             `json:"userId" openapi:"example:1"`
	User        *NotificationUser `json:"user" openapi:"$ref:NotificationUser"`
	Title       string            `json:"title" openapi:"example:title;required"`
	Description string            `json:"description" openapi:"example:description;required"`
	ReadAt      *time.Time        `json:"read_at" openapi:"example:2021-01-01T00:00:00Z"`
	ReadBy      *uint             `json:"-" openapi:"ignored"`
	ReadByUser  *NotificationUser `json:"readBy" openapi:"$ref:NotificationUser"`
	Extra       interface{}       `json:"-" openapi:"ignored"`
	IsForSystem bool              `json:"isForSystem" openapi:"example:false"`
	Status      *string           `json:"status" openapi:"example:pending"`
	StatusAt    *time.Time        `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt   time.Time         `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time         `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time        `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Notification) TableName() string {
	return "notifications"
}

func (u *Notification) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
