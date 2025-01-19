package domain

import "time"

/*
 * @apiDefine: CycleIncomingCyclePickupShiftStaff
 */
type CycleIncomingCyclePickupShiftStaff struct {
	ID        uint   `json:"id" openapi:"example:1"`
	UserID    uint   `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: CycleIncomingCyclePickupShiftCycleNextStaffType
 */
type CycleIncomingCyclePickupShiftCycleNextStaffType struct {
	ID        uint                    `json:"id" openapi:"example:1"`
	Role      *CycleNextStaffTypeRole `json:"role" openapi:"$ref:CycleNextStaffTypeRole"`
	DateTime  time.Time               `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName string                  `json:"shiftName" openapi:"example:morning"`
	StartHour time.Time               `json:"startHour" openapi:"example:00:00"`
	EndHour   time.Time               `json:"endHour" openapi:"example:00:00"`
}

/*
 * @apiDefine: CycleIncomingCyclePickupShift
 */
type CycleIncomingCyclePickupShift struct {
	ID                 uint                                             `json:"id" openapi:"example:1"`
	CycleID            uint                                             `json:"cycleId" openapi:"example:1"`
	Staff              *CycleIncomingCyclePickupShiftStaff              `json:"staff" openapi:"$ref:CycleIncomingCyclePickupShiftStaff"`
	CycleNextStaffType *CycleIncomingCyclePickupShiftCycleNextStaffType `json:"cycleNextStaffType" openapi:"$ref:CycleIncomingCyclePickupShiftCycleNextStaffType"`
	DateTime           time.Time                                        `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt          time.Time                                        `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          time.Time                                        `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *time.Time                                       `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}
