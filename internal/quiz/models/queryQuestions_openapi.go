package models

import "github.com/hoitek/Maja-Service/internal/quiz/domain"

/*
 * @apiDefine: QuizzesQueryQuestionsResponseData
 */
type QuizzesQueryQuestionsResponseData struct {
	Limit      int                                 `json:"limit" openapi:"example:10"`
	Offset     int                                 `json:"offset" openapi:"example:0"`
	Page       int                                 `json:"page" openapi:"example:1"`
	TotalRows  int                                 `json:"totalRows" openapi:"example:1"`
	TotalPages int                                 `json:"totalPages" openapi:"example:1"`
	Items      []QuizzesCreateQuestionResponseData `json:"items" openapi:"$ref:QuizzesCreateQuestionResponseData;type:array"`
}

/*
 * @apiDefine: QuizzesQueryQuestionsResponse
 */
type QuizzesQueryQuestionsResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       QuizzesQueryQuestionsResponseData `json:"data" openapi:"$ref:QuizzesQueryQuestionsResponseData"`
}

/*
 * @apiDefine: QuizzesQueryQuestionsNotFoundResponse
 */
type QuizzesQueryQuestionsNotFoundResponse struct {
	Quizzes []domain.Quiz `json:"quizzes" openapi:"$ref:Quiz;type:array"`
}
