package models

import "github.com/hoitek/Maja-Service/internal/quiz/domain"

/*
 * @apiDefine: QuizzesQueryResponseData
 */
type QuizzesQueryResponseData struct {
	Limit      int                         `json:"limit" openapi:"example:10"`
	Offset     int                         `json:"offset" openapi:"example:0"`
	Page       int                         `json:"page" openapi:"example:1"`
	TotalRows  int                         `json:"totalRows" openapi:"example:1"`
	TotalPages int                         `json:"totalPages" openapi:"example:1"`
	Items      []QuizzesCreateResponseData `json:"items" openapi:"$ref:QuizzesCreateResponseData;type:array"`
}

/*
 * @apiDefine: QuizzesQueryResponse
 */
type QuizzesQueryResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200"`
	Data       QuizzesQueryResponseData `json:"data" openapi:"$ref:QuizzesQueryResponseData"`
}

/*
 * @apiDefine: QuizzesQueryNotFoundResponse
 */
type QuizzesQueryNotFoundResponse struct {
	Quizzes []domain.Quiz `json:"quizzes" openapi:"$ref:Quiz;type:array"`
}
