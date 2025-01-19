package models

/*
 * @apiDefine: CyclesCreateVisitPauseResponse
 */
type CyclesCreateVisitPauseResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreatePickupShiftResponseData `json:"data" openapi:"$ref:CyclesCreatePickupShiftResponseData"`
}
