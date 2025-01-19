package domain

import (
	csDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	"time"
)

/*
 * @apiDefine: CyclePickupShiftStaffRolePermission
 */
type CyclePickupShiftStaffRolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:Dashboard"`
	Title string `json:"title" openapi:"example:Can See Dashboard"`
}

/*
 * @apiDefine: CyclePickupShiftStaffRole
 */
type CyclePickupShiftStaffRole struct {
	ID          uint                                  `json:"id" openapi:"example:1"`
	Name        string                                `json:"name" openapi:"example:staff;required"`
	Permissions []CyclePickupShiftStaffRolePermission `json:"permissions" openapi:"$ref:CyclePickupShiftStaffRolePermission;type:array"`
}

/*
 * @apiDefine: CyclePickupShiftStaff
 */
type CyclePickupShiftStaff struct {
	ID                          uint                         `json:"id" openapi:"example:1"`
	UserID                      uint                         `json:"userId" openapi:"example:1"`
	FirstName                   string                       `json:"firstName" openapi:"example:John;required"`
	LastName                    string                       `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl                   string                       `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
	Roles                       []*CyclePickupShiftStaffRole `json:"roles" openapi:"$ref:CyclePickupShiftStaffRole;type:array"`
	VehicleTypes                interface{}                  `json:"vehicleTypes" openapi:"example:[\"car\",\"bicycle\",\"public_transportation\"];type:array;"`
	VehicleLicenseTypes         interface{}                  `json:"vehicleLicenseTypes" openapi:"example:[\"automatic\",\"manual\"];type:array;"`
	SelectedVehicleTypeForCycle string                       `json:"selectedVehicleTypeForCycle" openapi:"example:car"`
	SelectedVehicleForCycle     string                       `json:"selectedVehicleForCycle" openapi:"example:own"`
}

/*
 * @apiDefine: CyclePickupShiftCycleStaffType
 */
type CyclePickupShiftCycleStaffType struct {
	ID          uint                `json:"id" openapi:"example:1"`
	Role        *CycleStaffTypeRole `json:"role" openapi:"$ref:CycleStaffTypeRole"`
	DateTime    time.Time           `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName   string              `json:"shiftName" openapi:"example:morning"`
	StartHour   time.Time           `json:"startHour" openapi:"example:00:00"`
	EndHour     time.Time           `json:"endHour" openapi:"example:00:00"`
	IsUnplanned bool                `json:"isUnplanned" openapi:"example:false"`
}

/*
 * @apiDefine: CyclePickupShift
 */
type CyclePickupShift struct {
	ID                      uint                            `json:"id" openapi:"example:1"`
	CycleID                 uint                            `json:"cycleId" openapi:"example:1"`
	Staff                   *CyclePickupShiftStaff          `json:"staff" openapi:"$ref:CyclePickupShiftStaff"`
	Shift                   *CycleShift                     `json:"shift" openapi:"$ref:CycleShift"`
	CycleStaffType          *CyclePickupShiftCycleStaffType `json:"cycleStaffType" openapi:"$ref:CyclePickupShiftCycleStaffType"`
	CustomerServices        []*csDomain.CustomerServices    `json:"customerServices" openapi:"$ref:CustomerServices"`
	Status                  string                          `json:"status" openapi:"example:not-started"`
	PrevStatus              string                          `json:"prevStatus" openapi:"example:not-started"`
	StartKilometer          *string                         `json:"startKilometer" openapi:"example:0"`
	ReasonOfTheCancellation *string                         `json:"reasonOfTheCancellation" openapi:"example:reason"`
	ReasonOfTheReactivation *string                         `json:"reasonOfTheReactivation" openapi:"example:reason"`
	ReasonOfTheResume       *string                         `json:"reasonOfTheResume" openapi:"example:reason"`
	ReasonOfThePause        *string                         `json:"reasonOfThePause" openapi:"example:reason"`
	IsUnplanned             bool                            `json:"isUnplanned" openapi:"example:false"`
	DateTime                time.Time                       `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt               time.Time                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt               time.Time                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt               *time.Time                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	StartedAt               *time.Time                      `json:"started_at" openapi:"example:2021-01-01T00:00:00Z"`
	EndedAt                 *time.Time                      `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CancelledAt             *time.Time                      `json:"cancelled_at" openapi:"example:2021-01-01T00:00:00Z"`
	DelayedAt               *time.Time                      `json:"delayed_at" openapi:"example:2021-01-01T00:00:00Z"`
	PausedAt                *time.Time                      `json:"paused_at" openapi:"example:2021-01-01T00:00:00Z"`
	ResumedAt               *time.Time                      `json:"resumed_at" openapi:"example:2021-01-01T00:00:00Z"`
	ReactivatedAt           *time.Time                      `json:"reactivated_at" openapi:"example:2021-01-01T00:00:00Z"`
}
