package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: ServiceGrade
 */
type ServiceGrade struct {
	ID          uint       `json:"id" openapi:"example:1"`
	Name        string     `json:"name" openapi:"example:John;required"`
	Description *string    `json:"description" openapi:"example:John;required"`
	Grade       int        `json:"grade" openapi:"example:0;required"`
	Color       string     `json:"color" openapi:"example:#000000;required"`
	CreatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *ServiceGrade) TableName() string {
	return "servicegrades"
}

func (u *ServiceGrade) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
