package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryVisitsTodosResponseData
 */
type CyclesQueryVisitsTodosResponseData struct {
	Limit      int                                 `json:"limit" openapi:"example:10"`
	Offset     int                                 `json:"offset" openapi:"example:0"`
	Page       int                                 `json:"page" openapi:"example:1"`
	TotalRows  int                                 `json:"totalRows" openapi:"example:1"`
	TotalPages int                                 `json:"totalPages" openapi:"example:1"`
	Items      []CyclesCreateVisitTodoResponseData `json:"items" openapi:"$ref:CyclesCreateVisitTodoResponseData;type:array"`
}

/*
 * @apiDefine: CyclesQueryVisitsTodosResponse
 */
type CyclesQueryVisitsTodosResponse struct {
	StatusCode int                                `json:"statusCode" openapi:"example:200;"`
	Date       CyclesQueryVisitsTodosResponseData `json:"data" openapi:"$ref:CyclesQueryVisitsTodosResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryVisitsTodosNotFoundResponse
 */
type CyclesQueryVisitsTodosNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
