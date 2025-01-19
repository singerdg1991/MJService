package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: Evaluation
 */
type Evaluation struct {
	ID             uint       `json:"id" openapi:"example:1"`
	StaffID        uint       `json:"staffId" openapi:"example:1;required"`
	EvaluationType string     `json:"evaluationType" openapi:"example:grace;required"`
	Title          string     `json:"title" openapi:"example:John;required"`
	Description    *string    `json:"description" openapi:"example:John;required"`
	CreatedAt      time.Time  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Evaluation) TableName() string {
	return "evaluations"
}

func (u *Evaluation) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
