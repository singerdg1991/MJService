package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryResponseDataItemStaffType
 */
type CyclesQueryResponseDataItemStaffType struct {
	ID               uint                      `json:"id" openapi:"example:1"`
	CycleID          uint                      `json:"cycleId" openapi:"example:1"`
	RoleID           uint                      `json:"roleId" openapi:"example:1"`
	Role             domain.CycleStaffTypeRole `json:"role" openapi:"$ref:CycleStaffTypeRole"`
	DateTime         string                    `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName        string                    `json:"shiftName" openapi:"example:morning"`
	NeededStaffCount uint                      `json:"neededStaffCount" openapi:"example:1"`
	StartHour        string                    `json:"startHour" openapi:"example:00:00"`
	EndHour          string                    `json:"endHour" openapi:"example:00:00"`
	UsedStaffCount   uint                      `json:"usedStaffCount" openapi:"example:1"`
	RemindStaffCount uint                      `json:"remindStaffCount" openapi:"example:1"`
}

/*
 * @apiDefine: CyclesQueryResponseDataItemNextStaffType
 */
type CyclesQueryResponseDataItemNextStaffType struct {
	ID               uint                      `json:"id" openapi:"example:1"`
	CurrentCycleID   uint                      `json:"currentCycleId" openapi:"example:1"`
	RoleID           uint                      `json:"roleId" openapi:"example:1"`
	Role             domain.CycleStaffTypeRole `json:"role" openapi:"$ref:CycleStaffTypeRole"`
	DateTime         string                    `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName        string                    `json:"shiftName" openapi:"example:morning"`
	NeededStaffCount uint                      `json:"neededStaffCount" openapi:"example:1"`
	StartHour        string                    `json:"startHour" openapi:"example:00:00"`
	EndHour          string                    `json:"endHour" openapi:"example:00:00"`
	UsedStaffCount   uint                      `json:"usedStaffCount" openapi:"example:1"`
	RemindStaffCount uint                      `json:"remindStaffCount" openapi:"example:1"`
}

/*
 * @apiDefine: CyclesQueryResponseDataItem
 */
type CyclesQueryResponseDataItem struct {
	ID                    uint                                       `json:"id" openapi:"example:1"`
	SectionID             uint                                       `json:"sectionId" openapi:"example:1;required;"`
	Name                  string                                     `json:"name" openapi:"example:John;required"`
	StartDate             string                                     `json:"start_date" openapi:"example:2021-01-01;required"`
	EndDate               string                                     `json:"end_date" openapi:"example:2021-01-01;required"`
	PeriodLength          string                                     `json:"periodLength" openapi:"example:oneWeek;required"`
	ShiftMorningStartTime string                                     `json:"shiftMorningStartTime" openapi:"example:08:00;required;"`
	ShiftMorningEndTime   string                                     `json:"shiftMorningEndTime" openapi:"example:16:00;required;"`
	ShiftEveningStartTime string                                     `json:"shiftEveningStartTime" openapi:"example:16:00;required;"`
	ShiftEveningEndTime   string                                     `json:"shiftEveningEndTime" openapi:"example:00:00;required;"`
	ShiftNightStartTime   string                                     `json:"shiftNightStartTime" openapi:"example:00:00;required;"`
	ShiftNightEndTime     string                                     `json:"shiftNightEndTime" openapi:"example:08:00;required;"`
	FreezePeriodDate      string                                     `json:"freeze_period_date" openapi:"example:2021-01-01;required"`
	WishDays              int                                        `json:"wishDays" openapi:"example:1;required"`
	StaffTypes            []CyclesQueryResponseDataItemStaffType     `json:"staffTypes" openapi:"$ref:CyclesQueryResponseDataItemStaffType;type:array"`
	NextStaffTypes        []CyclesQueryResponseDataItemNextStaffType `json:"nextStaffTypes" openapi:"$ref:CyclesQueryResponseDataItemNextStaffType;type:array"`
	Status                string                                     `json:"status" openapi:"example:active;required"`
	CreatedAt             string                                     `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             string                                     `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             string                                     `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryResponseData
 */
type CyclesQueryResponseData struct {
	Limit      int                           `json:"limit" openapi:"example:10"`
	Offset     int                           `json:"offset" openapi:"example:0"`
	Page       int                           `json:"page" openapi:"example:1"`
	TotalRows  int                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                           `json:"totalPages" openapi:"example:1"`
	Items      []CyclesQueryResponseDataItem `json:"items" openapi:"$ref:CyclesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesQueryResponse
 */
type CyclesQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryResponseData `json:"data" openapi:"$ref:CyclesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryNotFoundResponse
 */
type CyclesQueryNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
