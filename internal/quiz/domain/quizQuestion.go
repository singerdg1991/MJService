package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: QuizQuestionQuiz
 */
type QuizQuestionQuiz struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
}

/*
 * @apiDefine: QuizQuestionOptionItem
 */
type QuizQuestionOptionItem struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
	Score int    `json:"score" openapi:"example:1;required"`
}

/*
 * @apiDefine: QuizQuestion
 */
type QuizQuestion struct {
	ID          uint                      `json:"id" openapi:"example:1"`
	QuizID      uint                      `json:"quizId" openapi:"example:1"`
	Quiz        *QuizQuestionQuiz         `json:"quiz" openapi:"$ref:QuizQuestionQuiz"`
	Options     []*QuizQuestionOptionItem `json:"options" openapi:"$ref:QuizQuestionOptionItem;type:array"`
	Title       string                    `json:"title" openapi:"example:title;required"`
	Description *string                   `json:"description" openapi:"example:description"`
	CreatedAt   time.Time                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *QuizQuestion) TableName() string {
	return "quizQuestions"
}

func (u *QuizQuestion) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
