package models

/*
 * @apiDefine: CyclesUpdateVisitTodoStatusResponse
 */
type CyclesUpdateVisitTodoStatusResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateVisitTodoResponseData `json:"data" openapi:"$ref:CyclesCreateVisitTodoResponseData"`
}
