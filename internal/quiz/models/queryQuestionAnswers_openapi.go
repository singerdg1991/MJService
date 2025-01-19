package models

import "github.com/hoitek/Maja-Service/internal/quiz/domain"

/*
 * @apiDefine: QuizzesQueryQuestionAnswersResponseData
 */
type QuizzesQueryQuestionAnswersResponseData struct {
	Limit      int                               `json:"limit" openapi:"example:10"`
	Offset     int                               `json:"offset" openapi:"example:0"`
	Page       int                               `json:"page" openapi:"example:1"`
	TotalRows  int                               `json:"totalRows" openapi:"example:1"`
	TotalPages int                               `json:"totalPages" openapi:"example:1"`
	Items      []QuizzesCreateAnswerResponseData `json:"items" openapi:"$ref:QuizzesCreateAnswerResponseData;type:array"`
}

/*
 * @apiDefine: QuizzesQueryQuestionAnswersResponse
 */
type QuizzesQueryQuestionAnswersResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       QuizzesQueryQuestionAnswersResponseData `json:"data" openapi:"$ref:QuizzesQueryQuestionAnswersResponseData"`
}

/*
 * @apiDefine: QuizzesQueryQuestionAnswersNotFoundResponse
 */
type QuizzesQueryQuestionAnswersNotFoundResponse struct {
	Quizzes []domain.Quiz `json:"quizzes" openapi:"$ref:Quiz;type:array"`
}
