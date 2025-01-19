package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryPickupShiftsResponseData
 */
type CyclesQueryPickupShiftsResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []CyclesCreatePickupShiftResponseData `json:"items" openapi:"$ref:CyclesCreatePickupShiftResponseData;type:array"`
}

/*
 * @apiDefine: CyclesQueryPickupShiftsResponse
 */
type CyclesQueryPickupShiftsResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryPickupShiftsResponseData `json:"data" openapi:"$ref:CyclesQueryPickupShiftsResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryPickupShiftsNotFoundResponse
 */
type CyclesQueryPickupShiftsNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
