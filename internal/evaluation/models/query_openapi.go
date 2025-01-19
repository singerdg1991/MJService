package models

import "github.com/hoitek/Maja-Service/internal/evaluation/domain"

/*
 * @apiDefine: EvaluationsQueryResponseData
 */
type EvaluationsQueryResponseData struct {
	Limit      int                       `json:"limit" openapi:"example:10"`
	Offset     int                       `json:"offset" openapi:"example:0"`
	Page       int                       `json:"page" openapi:"example:1"`
	TotalRows  int                       `json:"totalRows" openapi:"example:1"`
	TotalPages int                       `json:"totalPages" openapi:"example:1"`
	Items      []EvaluationsResponseData `json:"items" openapi:"$ref:EvaluationsResponseData;type:array"`
}

/*
 * @apiDefine: EvaluationsQueryResponse
 */
type EvaluationsQueryResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200;"`
	Data       EvaluationsQueryResponseData `json:"data" openapi:"$ref:EvaluationsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: EvaluationsQueryNotFoundResponse
 */
type EvaluationsQueryNotFoundResponse struct {
	Evaluations []domain.Evaluation `json:"evaluations" openapi:"$ref:Evaluation;type:array"`
}
