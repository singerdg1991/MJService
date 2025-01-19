package models

import (
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
)

/*
 * @apiDefine: QuizzesCreateQuestionResponseDataPermission
 */
type QuizzesCreateQuestionResponseDataPermission struct {
	ID    int64  `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:John;required"`
	Title string `json:"title" openapi:"example:John;required"`
}

/*
 * @apiDefine: QuizzesCreateQuestionResponseData
 */
type QuizzesCreateQuestionResponseData struct {
	ID          uint                             `json:"id" openapi:"example:1"`
	QuizID      uint                             `json:"quizId" openapi:"example:1;required"`
	Quiz        *domain.QuizQuestionQuiz         `json:"quiz" openapi:"$ref:QuizQuestionQuiz"`
	Options     []*domain.QuizQuestionOptionItem `json:"options" openapi:"$ref:QuizQuestionOptionItem;type:array"`
	Title       string                           `json:"title" openapi:"example:title;required"`
	Description *string                          `json:"description" openapi:"example:description;required"`
	CreatedAt   string                           `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   string                           `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *string                          `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesCreateQuestionResponse
 */
type QuizzesCreateQuestionResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       []QuizzesCreateQuestionResponseData `json:"data" openapi:"$ref:QuizzesCreateQuestionResponseData;type:array;required"`
}
