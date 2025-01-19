package models

import "github.com/hoitek/Maja-Service/internal/todo/domain"

/*
 * @apiDefine: TodosQueryResponseData
 */
type TodosQueryResponseData struct {
	Limit      int                 `json:"limit" openapi:"example:10"`
	Offset     int                 `json:"offset" openapi:"example:0"`
	Page       int                 `json:"page" openapi:"example:1"`
	TotalRows  int                 `json:"totalRows" openapi:"example:1"`
	TotalPages int                 `json:"totalPages" openapi:"example:1"`
	Items      []TodosResponseData `json:"items" openapi:"$ref:TodosResponseData;type:array"`
}

/*
 * @apiDefine: TodosQueryResponse
 */
type TodosQueryResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200;"`
	Data       TodosQueryResponseData `json:"data" openapi:"$ref:TodosQueryResponseData;type:object;"`
}

/*
 * @apiDefine: TodosQueryNotFoundResponse
 */
type TodosQueryNotFoundResponse struct {
	Todos []domain.Todo `json:"todos" openapi:"$ref:Todo;type:array"`
}
