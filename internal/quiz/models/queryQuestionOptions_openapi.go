package models

import "github.com/hoitek/Maja-Service/internal/quiz/domain"

/*
 * @apiDefine: QuizzesQueryQuestionOptionsResponseDataItem
 */
type QuizzesQueryQuestionOptionsResponseDataItem struct {
	ID             uint    `json:"id" openapi:"example:1"`
	QuizQuestionId uint    `json:"quizQuestionId" openapi:"example:1"`
	Title          string  `json:"title" openapi:"example:title"`
	Score          int     `json:"score" openapi:"example:1"`
	CreatedAt      string  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      string  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *string `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsResponseData
 */
type QuizzesQueryQuestionOptionsResponseData struct {
	Limit      int                                           `json:"limit" openapi:"example:10"`
	Offset     int                                           `json:"offset" openapi:"example:0"`
	Page       int                                           `json:"page" openapi:"example:1"`
	TotalRows  int                                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                                           `json:"totalPages" openapi:"example:1"`
	Items      []QuizzesQueryQuestionOptionsResponseDataItem `json:"items" openapi:"$ref:QuizzesQueryQuestionOptionsResponseDataItem;type:array"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsResponse
 */
type QuizzesQueryQuestionOptionsResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       QuizzesQueryQuestionOptionsResponseData `json:"data" openapi:"$ref:QuizzesQueryQuestionOptionsResponseData"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsNotFoundResponse
 */
type QuizzesQueryQuestionOptionsNotFoundResponse struct {
	Quizzes []domain.Quiz `json:"quizzes" openapi:"$ref:Quiz;type:array"`
}
