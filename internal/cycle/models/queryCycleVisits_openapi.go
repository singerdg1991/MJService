package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemStaffType
 */
type CyclesQueryCycleVisitsResponseDataItemStaffType struct {
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
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemNextStaffType
 */
type CyclesQueryCycleVisitsResponseDataItemNextStaffType struct {
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
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemVisitShiftStaffRole
 */
type CyclesQueryCycleVisitsResponseDataItemVisitShiftStaffRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:Doctor"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemVisitShiftStaff
 */
type CyclesQueryCycleVisitsResponseDataItemVisitShiftStaff struct {
	ID                    uint                                                      `json:"id" openapi:"example:1"`
	FirstName             string                                                    `json:"firstName" openapi:"example:John"`
	LastName              string                                                    `json:"lastName" openapi:"example:Doe"`
	AvatarUrl             string                                                    `json:"avatarUrl" openapi:"example:https://example.com/avatar.png"`
	Role                  CyclesQueryCycleVisitsResponseDataItemVisitShiftStaffRole `json:"role" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemVisitShiftStaffRole"`
	VisitsDoneCount       uint                                                      `json:"visitsDoneCount" openapi:"example:1"`
	VisitsCancelledCount  uint                                                      `json:"visitsCancelledCount" openapi:"example:1"`
	VisitsDelayCount      uint                                                      `json:"visitsDelayCount" openapi:"example:1"`
	VisitsNotStartedCount uint                                                      `json:"visitsNotStartedCount" openapi:"example:1"`
	VisitsOnGoingCount    uint                                                      `json:"visitsOnGoingCount" openapi:"example:1"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer
 */
type CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.png"`
	Gender    string `json:"gender" openapi:"example:male"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisitPaySoonBy
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisitPaySoonBy struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John"`
	LastName  string `json:"lastName" openapi:"example:Doe"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.png"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisitKey
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisitKey struct {
	ID  uint   `json:"id" openapi:"example:1"`
	Key string `json:"key" openapi:"example:13456"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressStreet
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressStreet struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:Street 1"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressCity
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressCity struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:City 1"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisitAddress
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisitAddress struct {
	ID     uint                                                          `json:"id" openapi:"example:1"`
	Name   string                                                        `json:"name" openapi:"example:John"`
	Street CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressStreet `json:"street" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressStreet"`
	City   CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressCity   `json:"city" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisitAddressCity"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShiftVisit
 */
type CyclesQueryCycleVisitsResponseDataItemShiftVisit struct {
	ID             uint                                                            `json:"id" openapi:"example:1"`
	Customer       CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer `json:"customer" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer"`
	ShiftName      string                                                          `json:"shiftName" openapi:"example:morning"`
	VisitStartTime string                                                          `json:"visitStartTime" openapi:"example:2021-01-01T00:00:00Z"`
	VisitEndTime   string                                                          `json:"visitEndTime" openapi:"example:2021-01-01T00:00:00Z"`
	VisitStatus    string                                                          `json:"visitStatus" openapi:"example:ongoing"`
	PaySoonAt      string                                                          `json:"pay_soon_at" openapi:"example:2021-01-01T00:00:00Z"`
	PaySoonBy      CyclesQueryCycleVisitsResponseDataItemShiftVisitPaySoonBy       `json:"pay_soon_by" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisitPaySoonBy"`
	VisitLength    string                                                          `json:"visitLength" openapi:"example:30min"`
	Keys           []CyclesQueryCycleVisitsResponseDataItemShiftVisitKey           `json:"keys" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisitKey"`
	Address        CyclesQueryCycleVisitsResponseDataItemShiftVisitAddress         `json:"address" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisitAddress"`
	Date           string                                                          `json:"date" openapi:"example:2021-01-01"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseDataItemShift
 */
type CyclesQueryCycleVisitsResponseDataItemShift struct {
	Staff           []CyclesQueryCycleVisitsResponseDataItemVisitShiftStaff         `json:"staff" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemVisitShiftStaff"`
	OnGoingCustomer CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer `json:"onGoingCustomer" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemVisitShiftOnGoingCustomer"`
	Visits          []CyclesQueryCycleVisitsResponseDataItemShiftVisit              `json:"visits" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShiftVisit"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponseData
 */
type CyclesQueryCycleVisitsResponseData struct {
	ID                    uint                                                  `json:"id" openapi:"example:1"`
	SectionID             int                                                   `json:"sectionId" openapi:"example:1;required;"`
	Name                  string                                                `json:"name" openapi:"example:John;required"`
	StartDate             string                                                `json:"start_date" openapi:"example:2021-01-01;required"`
	EndDate               string                                                `json:"end_date" openapi:"example:2021-01-01;required"`
	PeriodLength          string                                                `json:"periodLength" openapi:"example:oneWeek;required"`
	ShiftMorningStartTime string                                                `json:"shiftMorningStartTime" openapi:"example:08:00;required;"`
	ShiftMorningEndTime   string                                                `json:"shiftMorningEndTime" openapi:"example:16:00;required;"`
	ShiftEveningStartTime string                                                `json:"shiftEveningStartTime" openapi:"example:16:00;required;"`
	ShiftEveningEndTime   string                                                `json:"shiftEveningEndTime" openapi:"example:00:00;required;"`
	ShiftNightStartTime   string                                                `json:"shiftNightStartTime" openapi:"example:00:00;required;"`
	ShiftNightEndTime     string                                                `json:"shiftNightEndTime" openapi:"example:08:00;required;"`
	FreezePeriodDate      string                                                `json:"freeze_period_date" openapi:"example:2021-01-01;required"`
	WishDays              int                                                   `json:"wishDays" openapi:"example:1;required"`
	StaffTypes            []CyclesQueryCycleVisitsResponseDataItemStaffType     `json:"staffTypes" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemStaffType;type:array"`
	NextStaffTypes        []CyclesQueryCycleVisitsResponseDataItemNextStaffType `json:"nextStaffTypes" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemNextStaffType;type:array"`
	Shifts                []CyclesQueryCycleVisitsResponseDataItemShift         `json:"shifts" openapi:"$ref:CyclesQueryCycleVisitsResponseDataItemShift;type:array"`
	Status                string                                                `json:"status" openapi:"example:active;required"`
	CreatedAt             string                                                `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             string                                                `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             string                                                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsResponse
 */
type CyclesQueryCycleVisitsResponse struct {
	StatusCode int                                `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryCycleVisitsResponseData `json:"data" openapi:"$ref:CyclesQueryCycleVisitsResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryCycleVisitsNotFoundResponse
 */
type CyclesQueryCycleVisitsNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
