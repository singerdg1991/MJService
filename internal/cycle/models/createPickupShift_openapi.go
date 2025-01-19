package models

import (
	csDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreatePickupShiftResponseData
 */
type CyclesCreatePickupShiftResponseData struct {
	ID                      uint                                   `json:"id" openapi:"example:1"`
	CycleID                 uint                                   `json:"cycleId" openapi:"example:1"`
	Staff                   *domain.CyclePickupShiftStaff          `json:"staff" openapi:"$ref:CyclePickupShiftStaff"`
	Shift                   *domain.CycleShift                     `json:"shift" openapi:"$ref:CycleShift"`
	CycleStaffType          *domain.CyclePickupShiftCycleStaffType `json:"cycleStaffType" openapi:"$ref:CyclePickupShiftCycleStaffType"`
	CustomerServices        []*csDomain.CustomerServices           `json:"customerServices" openapi:"$ref:CustomerServices"`
	DateTime                string                                 `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	Status                  string                                 `json:"status" openapi:"example:not-started"`
	PrevStatus              string                                 `json:"prevStatus" openapi:"example:not-started"`
	StartKilometer          *string                                `json:"startKilometer" openapi:"example:0"`
	ReasonOfTheCancellation *string                                `json:"reasonOfTheCancellation" openapi:"example:reason"`
	ReasonOfTheReactivation *string                                `json:"reasonOfTheReactivation" openapi:"example:reason"`
	ReasonOfTheResume       *string                                `json:"reasonOfTheResume" openapi:"example:reason"`
	ReasonOfThePause        *string                                `json:"reasonOfThePause" openapi:"example:reason"`
	IsUnplanned             bool                                   `json:"isUnplanned" openapi:"example:false"`
	CreatedAt               string                                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt               string                                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt               *string                                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	StartedAt               *string                                `json:"started_at" openapi:"example:2021-01-01T00:00:00Z"`
	EndedAt                 *string                                `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CancelledAt             *string                                `json:"cancelled_at" openapi:"example:2021-01-01T00:00:00Z"`
	DelayedAt               *string                                `json:"delayed_at" openapi:"example:2021-01-01T00:00:00Z"`
	PausedAt                *string                                `json:"paused_at" openapi:"example:2021-01-01T00:00:00Z"`
	ResumedAt               *string                                `json:"resumed_at" openapi:"example:2021-01-01T00:00:00Z"`
	ReactivatedAt           *string                                `json:"reactivated_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreatePickupShiftResponse
 */
type CyclesCreatePickupShiftResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesQueryPickupShiftsResponseData `json:"data" openapi:"$ref:CyclesQueryPickupShiftsResponseData"`
}
