package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateIncomingShiftSwapResponseData
 */
type CyclesCreateIncomingShiftSwapResponseData struct {
	SourceShifts []*domain.CycleIncomingCyclePickupShift `json:"sourceShift" openapi:"$ref:CycleIncomingCyclePickupShift;type:array"`
	TargetShifts []*domain.CycleIncomingCyclePickupShift `json:"targetShift" openapi:"$ref:CycleIncomingCyclePickupShift;type:array"`
}

/*
 * @apiDefine: CyclesCreateIncomingShiftSwapResponse
 */
type CyclesCreateIncomingShiftSwapResponse struct {
	StatusCode int                                       `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateIncomingShiftSwapResponseData `json:"data" openapi:"$ref:CyclesCreateIncomingShiftSwapResponseData"`
}
