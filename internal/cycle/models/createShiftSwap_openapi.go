package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateShiftSwapResponseData
 */
type CyclesCreateShiftSwapResponseData struct {
	SourceShifts []*domain.CyclePickupShift `json:"sourceShift" openapi:"$ref:CyclePickupShift;type:array"`
	TargetShifts []*domain.CyclePickupShift `json:"targetShift" openapi:"$ref:CyclePickupShift;type:array"`
}

/*
 * @apiDefine: CyclesCreateShiftSwapResponse
 */
type CyclesCreateShiftSwapResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateShiftSwapResponseData `json:"data" openapi:"$ref:CyclesCreateShiftSwapResponseData"`
}
