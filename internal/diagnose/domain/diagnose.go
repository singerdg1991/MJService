package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: Diagnose
 */
type Diagnose struct {
	ID          uint       `json:"id" openapi:"example:1"`
	Title       string     `json:"title" openapi:"example:John;required"`
	Code        string     `json:"code" openapi:"example:John;required"`
	Description *string    `json:"description" openapi:"example:John;required"`
	CreatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Diagnose) TableName() string {
	return "diagnoses"
}

func (u *Diagnose) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
