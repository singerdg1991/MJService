package domain

import "time"

/*
 * @apiDefine: CycleShiftCustomerHomeKeyRole
 */
type CycleShiftCustomerHomeKeyRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: CycleShiftCustomerHomeKeyCreatedBy
 */
type CycleShiftCustomerHomeKeyCreatedBy struct {
	ID        int64  `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: CycleShiftCustomerHomeKey
 */
type CycleShiftCustomerHomeKey struct {
	ID        uint                                `json:"id" openapi:"example:1"`
	ShiftID   uint                                `json:"shiftId" openapi:"example:1"`
	KeyNo     string                              `json:"keyNo" openapi:"example:1"`
	Status    string                              `json:"status" openapi:"example:accepted"`
	Reason    *string                             `json:"reason" openapi:"example:accepted"`
	CreatedAt time.Time                           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy *CycleShiftCustomerHomeKeyCreatedBy `json:"createdBy" openapi:"$ref:CycleShiftCustomerHomeKeyCreatedBy;required"`
	UpdatedAt time.Time                           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *time.Time                          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}
