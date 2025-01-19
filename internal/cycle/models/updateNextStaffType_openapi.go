package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesUpdateNextStaffTypeResponseDataItem
 */
type CyclesUpdateNextStaffTypeResponseDataItem struct {
	ID               uint                      `json:"id" openapi:"example:1"`
	CurrentCycleID   uint                      `json:"currentCycleId" openapi:"example:1"`
	RoleID           uint                      `json:"roleId" openapi:"example:1"`
	Role             domain.CycleStaffTypeRole `json:"role" openapi:"$ref:CycleStaffTypeRole"`
	DateTime         string                    `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName        string                    `json:"shiftName" openapi:"example:morning"`
	StartHour        string                    `json:"startHour" openapi:"example:00:00"`
	EndHour          string                    `json:"endHour" openapi:"example:00:00"`
	NeededStaffCount uint                      `json:"neededStaffCount" openapi:"example:1"`
	UsedStaffCount   uint                      `json:"usedStaffCount" openapi:"example:1"`
	RemindStaffCount uint                      `json:"remindStaffCount" openapi:"example:1"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypeResponseData
 */
type CyclesUpdateNextStaffTypeResponseData struct {
	Limit      int                                         `json:"limit" openapi:"example:10"`
	Offset     int                                         `json:"offset" openapi:"example:0"`
	Page       int                                         `json:"page" openapi:"example:1"`
	TotalRows  int                                         `json:"totalRows" openapi:"example:1"`
	TotalPages int                                         `json:"totalPages" openapi:"example:1"`
	Items      []CyclesUpdateNextStaffTypeResponseDataItem `json:"items" openapi:"$ref:CyclesUpdateNextStaffTypeResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypeResponse
 */
type CyclesUpdateNextStaffTypeResponse struct {
	StatusCode int                                   `json:"statusCode" openapi:"example:200;"`
	Date       CyclesUpdateNextStaffTypeResponseData `json:"data" openapi:"$ref:CyclesUpdateNextStaffTypeResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypeNotFoundResponse
 */
type CyclesUpdateNextStaffTypeNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
