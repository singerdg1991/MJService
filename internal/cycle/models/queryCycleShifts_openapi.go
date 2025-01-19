package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryCycleShiftsResponseDataItem
 */
type CyclesQueryCycleShiftsResponseDataItem struct {
	ID            uint                     `json:"id" openapi:"example:1"`
	ExchangeKey   string                   `json:"exchangeKey" openapi:"example:dfhdsjrtwerwrwfgjgfrt"`
	CycleID       uint                     `json:"cycleId" openapi:"example:1"`
	StaffTypeIDs  []uint                   `json:"staffTypeIds" openapi:"example:[1,2,3]"`
	StaffTypes    []*domain.CycleStaffType `json:"staffTypes" openapi:"$ref:CycleStaffType;type:array"`
	ShiftName     string                   `json:"shiftName" openapi:"example:morning"`
	VehicleType   *string                  `json:"vehicleType" openapi:"example:own"`
	StartLocation *string                  `json:"startLocation" openapi:"example:office"`
	DateTime      string                   `json:"dateTime" openapi:"example:2021-08-02"`
	Status        string                   `json:"status" openapi:"example:not-started"`
	CreatedAt     string                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     string                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *string                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryCycleShiftsResponseData
 */
type CyclesQueryCycleShiftsResponseData struct {
	Limit      int                                      `json:"limit" openapi:"example:10"`
	Offset     int                                      `json:"offset" openapi:"example:0"`
	Page       int                                      `json:"page" openapi:"example:1"`
	TotalRows  int                                      `json:"totalRows" openapi:"example:1"`
	TotalPages int                                      `json:"totalPages" openapi:"example:1"`
	Items      []CyclesQueryCycleShiftsResponseDataItem `json:"items" openapi:"$ref:CyclesQueryCycleShiftsResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesQueryCycleShiftsResponse
 */
type CyclesQueryCycleShiftsResponse struct {
	StatusCode int                                `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryCycleShiftsResponseData `json:"data" openapi:"$ref:CyclesQueryCycleShiftsResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryCycleShiftsNotFoundResponse
 */
type CyclesQueryCycleShiftsNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
