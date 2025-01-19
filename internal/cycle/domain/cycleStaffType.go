package domain

import "time"

/*
 * @apiDefine: CycleStaffTypeRole
 */
type CycleStaffTypeRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: CycleStaffType
 */
type CycleStaffType struct {
	ID               uint                `json:"id" openapi:"example:1"`
	CycleID          uint                `json:"cycleId" openapi:"example:1"`
	RoleID           uint                `json:"roleId" openapi:"example:1"`
	Role             *CycleStaffTypeRole `json:"role" openapi:"$ref:CycleStaffTypeRole"`
	DateTime         time.Time           `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	ShiftName        string              `json:"shiftName" openapi:"example:morning"`
	NeededStaffCount uint                `json:"neededStaffCount" openapi:"example:1"`
	StartHour        time.Time           `json:"startHour" openapi:"example:00:00"`
	EndHour          time.Time           `json:"endHour" openapi:"example:00:00"`
	IsUnplanned      bool                `json:"isUnplanned" openapi:"example:false"`
	UsedStaffCount   uint                `json:"usedStaffCount" openapi:"example:1"`
	RemindStaffCount uint                `json:"remindStaffCount" openapi:"example:1"`
	CreatedAt        time.Time           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt        time.Time           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt        *time.Time          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}
