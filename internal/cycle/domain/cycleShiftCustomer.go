package domain

import "time"

/*
 * @apiDefine: CycleShiftCustomerCustomerUser
 */
type CycleShiftCustomerCustomerUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: CycleShiftCustomerCustomer
 */
type CycleShiftCustomerCustomer struct {
	ID     uint                            `json:"id" openapi:"example:1"`
	UserID uint                            `json:"userId" openapi:"example:1"`
	User   *CycleShiftCustomerCustomerUser `json:"user" openapi:"$ref:CycleShiftCustomerCustomerUser"`
	KeyNo  string                          `json:"keyNo" openapi:"example:123456"`
}

/*
 * @apiDefine: CycleShiftCustomer
 */
type CycleShiftCustomer struct {
	ID          uint                        `json:"id" openapi:"example:1"`
	CycleID     uint                        `json:"cycleId" openapi:"example:1"`
	StaffTypeID uint                        `json:"staffTypeId" openapi:"example:1"`
	CustomerID  uint                        `json:"customerId" openapi:"example:1"`
	Customer    *CycleShiftCustomerCustomer `json:"customer" openapi:"$ref:CycleShiftCustomerCustomer"`
	CreatedAt   time.Time                   `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                   `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}
