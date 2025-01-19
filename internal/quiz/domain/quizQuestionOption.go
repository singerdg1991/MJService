package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: QuizQuestionOption
 */
type QuizQuestionOption struct {
	ID             uint       `json:"id" openapi:"example:1"`
	QuizQuestionId uint       `json:"quizQuestionId" openapi:"example:1"`
	Title          string     `json:"title" openapi:"example:title;required"`
	Score          int        `json:"score" openapi:"example:0"`
	CreatedAt      time.Time  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *QuizQuestionOption) TableName() string {
	return "quizQuestionOptions"
}

func (u *QuizQuestionOption) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
