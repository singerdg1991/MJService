package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: Push
 */
type Push struct {
	ID         uint       `json:"id" openapi:"example:1"`
	UserID     uint       `json:"userId" openapi:"example:1"`
	Endpoint   string     `json:"endpoint" openapi:"example:endpoint;required"`
	KeysAuth   string     `json:"keysAuth" openapi:"example:keysAuth;required"`
	KeysP256dh string     `json:"keysP256dh" openapi:"example:keysP256dh;required"`
	CreatedAt  time.Time  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt  time.Time  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt  *time.Time `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Push) TableName() string {
	return "pushes"
}

func (u *Push) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
