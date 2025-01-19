package models

import "github.com/hoitek/Maja-Service/internal/todo/domain"

/*
 * @apiDefine: TodosResponseData
 */
type TodosResponseData struct {
	ID          uint            `json:"id" openapi:"example:1"`
	UserID      uint            `json:"userId" openapi:"example:1"`
	Title       string          `json:"title" openapi:"example:title"`
	Date        string          `json:"date" openapi:"example:2021-01-01"`
	Time        string          `json:"time" openapi:"example:00:00"`
	User        domain.TodoUser `json:"user" openapi:"$ref:TodoUser"`
	Description *string         `json:"description" openapi:"example:test"`
	Status      string          `json:"status" openapi:"example:active"`
	DoneAt      string          `json:"done_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy   domain.TodoUser `json:"createdBy" openapi:"$ref:TodoUser"`
}

/*
 * @apiDefine: TodosCreateResponse
 */
type TodosCreateResponse struct {
	StatusCode int               `json:"statusCode" openapi:"example:200"`
	Data       TodosResponseData `json:"data" openapi:"$ref:TodosResponseData"`
}
