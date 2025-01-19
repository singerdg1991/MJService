package models

/*
 * @apiDefine: CyclesCreateVisitDelayResponse
 */
type CyclesCreateVisitDelayResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreatePickupShiftResponseData `json:"data" openapi:"$ref:CyclesCreatePickupShiftResponseData"`
}
