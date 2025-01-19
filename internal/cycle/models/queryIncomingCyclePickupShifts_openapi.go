package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryIncomingCyclePickupShiftsResponseData
 */
type CyclesQueryIncomingCyclePickupShiftsResponseData struct {
	Limit      int                                                `json:"limit" openapi:"example:10"`
	Offset     int                                                `json:"offset" openapi:"example:0"`
	Page       int                                                `json:"page" openapi:"example:1"`
	TotalRows  int                                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                                `json:"totalPages" openapi:"example:1"`
	Items      []CyclesCreateIncomingCyclePickupShiftResponseData `json:"items" openapi:"$ref:CyclesCreateIncomingCyclePickupShiftResponseData;type:array"`
}

/*
 * @apiDefine: CyclesQueryIncomingCyclePickupShiftsResponse
 */
type CyclesQueryIncomingCyclePickupShiftsResponse struct {
	StatusCode int                                              `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryIncomingCyclePickupShiftsResponseData `json:"data" openapi:"$ref:CyclesQueryIncomingCyclePickupShiftsResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryIncomingCyclePickupShiftsNotFoundResponse
 */
type CyclesQueryIncomingCyclePickupShiftsNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
