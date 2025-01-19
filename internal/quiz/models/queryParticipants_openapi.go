package models

import (
	"github.com/hoitek/Maja-Service/internal/quiz/domain"
)

/*
 * @apiDefine: QuizzesQueryParticipantsResponseDataItem
 */
type QuizzesQueryParticipantsResponseDataItem struct {
	ID        uint                        `json:"id" openapi:"example:1"`
	QuizID    uint                        `json:"quizId" openapi:"example:1"`
	UserID    uint                        `json:"userId" openapi:"example:1"`
	Quiz      *domain.QuizParticipantQuiz `json:"quiz" openapi:"$ref:QuizParticipantQuiz"`
	User      *domain.QuizParticipantUser `json:"user" openapi:"$ref:QuizParticipantUser"`
	EndedAt   *string                     `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt string                      `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt string                      `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *string                     `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesQueryParticipantsResponseData
 */
type QuizzesQueryParticipantsResponseData struct {
	Limit      int                                        `json:"limit" openapi:"example:10"`
	Offset     int                                        `json:"offset" openapi:"example:0"`
	Page       int                                        `json:"page" openapi:"example:1"`
	TotalRows  int                                        `json:"totalRows" openapi:"example:1"`
	TotalPages int                                        `json:"totalPages" openapi:"example:1"`
	Items      []QuizzesQueryParticipantsResponseDataItem `json:"items" openapi:"$ref:QuizzesQueryParticipantsResponseDataItem;type:array"`
}

/*
 * @apiDefine: QuizzesQueryParticipantsResponse
 */
type QuizzesQueryParticipantsResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       QuizzesQueryParticipantsResponseData `json:"data" openapi:"$ref:QuizzesQueryParticipantsResponseData"`
}

/*
 * @apiDefine: QuizzesQueryParticipantsNotFoundResponse
 */
type QuizzesQueryParticipantsNotFoundResponse struct {
	Quizzes []domain.Quiz `json:"quizzes" openapi:"$ref:Quiz;type:array"`
}
