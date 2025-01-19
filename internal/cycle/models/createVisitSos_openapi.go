package models

/*
 * @apiDefine: CyclesCreateVisitSosResponseData
 */
type CyclesCreateVisitSosResponseData struct {
	TimeValueInSeconds int `json:"timeValueInSeconds" openapi:"example:1"`
}

/*
 * @apiDefine: CyclesCreateVisitSosResponse
 */
type CyclesCreateVisitSosResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateVisitSosResponseData `json:"data" openapi:"$ref:CyclesCreateVisitSosResponseData"`
}
