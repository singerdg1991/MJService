package models

import (
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
)

/*
 * @apiDefine: QuizzesCreateStartResponseData
 */
type QuizzesCreateStartResponseData struct {
	ID        uint                        `json:"id" openapi:"example:1"`
	QuizID    uint                        `json:"quizId" openapi:"example:1;required"`
	UserID    uint                        `json:"userId" openapi:"example:1;required"`
	Quiz      *domain.QuizParticipantQuiz `json:"quiz" openapi:"$ref:QuizParticipantQuiz"`
	User      *domain.QuizParticipantUser `json:"user" openapi:"$ref:QuizParticipantUser"`
	CreatedAt string                      `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt string                      `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *string                     `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesCreateStartResponse
 */
type QuizzesCreateStartResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       []QuizzesCreateStartResponseData `json:"data" openapi:"$ref:QuizzesCreateStartResponseData;type:array;required"`
}
