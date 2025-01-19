package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryCurrentResponseDataItemStaffType
 */
type CyclesQueryCurrentResponseDataItemStaffType struct {
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
 * @apiDefine: CyclesQueryCurrentResponseDataItemNextStaffType
 */
type CyclesQueryCurrentResponseDataItemNextStaffType struct {
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
 * @apiDefine: CyclesQueryCurrentResponseData
 */
type CyclesQueryCurrentResponseData struct {
	ID                    uint                                              `json:"id" openapi:"example:1"`
	SectionID             uint                                              `json:"sectionId" openapi:"example:1;required;"`
	Name                  string                                            `json:"name" openapi:"example:John;required"`
	StartDate             string                                            `json:"start_date" openapi:"example:2021-01-01;required"`
	EndDate               string                                            `json:"end_date" openapi:"example:2021-01-01;required"`
	PeriodLength          string                                            `json:"periodLength" openapi:"example:oneWeek;required"`
	ShiftMorningStartTime string                                            `json:"shiftMorningStartTime" openapi:"example:08:00;required;"`
	ShiftMorningEndTime   string                                            `json:"shiftMorningEndTime" openapi:"example:16:00;required;"`
	ShiftEveningStartTime string                                            `json:"shiftEveningStartTime" openapi:"example:16:00;required;"`
	ShiftEveningEndTime   string                                            `json:"shiftEveningEndTime" openapi:"example:00:00;required;"`
	ShiftNightStartTime   string                                            `json:"shiftNightStartTime" openapi:"example:00:00;required;"`
	ShiftNightEndTime     string                                            `json:"shiftNightEndTime" openapi:"example:08:00;required;"`
	FreezePeriodDate      string                                            `json:"freeze_period_date" openapi:"example:2021-01-01;required"`
	WishDays              int                                               `json:"wishDays" openapi:"example:1;required"`
	StaffTypes            []CyclesQueryCurrentResponseDataItemStaffType     `json:"staffTypes" openapi:"$ref:CyclesQueryCurrentResponseDataItemStaffType;type:array"`
	NextStaffTypes        []CyclesQueryCurrentResponseDataItemNextStaffType `json:"nextStaffTypes" openapi:"$ref:CyclesQueryCurrentResponseDataItemNextStaffType;type:array" bson:"cycleNextStaffTypes"`
	Status                string                                            `json:"status" openapi:"example:active;required"`
	CreatedAt             string                                            `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             string                                            `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             string                                            `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryCurrentResponse
 */
type CyclesQueryCurrentResponse struct {
	StatusCode int                            `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryCurrentResponseData `json:"data" openapi:"$ref:CyclesQueryCurrentResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryCurrentNotFoundResponse
 */
type CyclesQueryCurrentNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
