package models

import (
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
)

/*
 * @apiDefine: QuizzesCreateAnswerResponseData
 */
type QuizzesCreateAnswerResponseData struct {
	ID                   uint                                         `json:"id" openapi:"example:1"`
	UserID               uint                                         `json:"userId" openapi:"example:1"`
	User                 *domain.QuizQuestionAnswerUser               `json:"user" openapi:"$ref:QuizQuestionAnswerUser"`
	QuestionID           uint                                         `json:"questionId" openapi:"example:1"`
	Question             *domain.QuizQuestionAnswerQuestion           `json:"question" openapi:"$ref:QuizQuestionAnswerQuestion"`
	QuizQuestionOptionId uint                                         `json:"quizQuestionOptionId" openapi:"example:1"`
	QuizQuestionOption   *domain.QuizQuestionAnswerQuizQuestionOption `json:"quizQuestionOption" openapi:"$ref:QuizQuestionAnswerQuizQuestionOption"`
	CreatedAt            string                                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt            string                                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt            *string                                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesCreateAnswerResponse
 */
type QuizzesCreateAnswerResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       []QuizzesCreateAnswerResponseData `json:"data" openapi:"$ref:QuizzesCreateAnswerResponseData;type:array;required"`
}
