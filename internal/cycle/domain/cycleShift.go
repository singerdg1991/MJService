package domain

import "time"

/*
 * @apiDefine: CycleShift
 */
type CycleShift struct {
	ID            uint              `json:"id" openapi:"example:1"`
	ExchangeKey   string            `json:"exchangeKey" openapi:"example:dfhdsjrtwerwrwfgjgfrt"`
	CycleID       uint              `json:"cycleId" openapi:"example:1"`
	StaffTypeIDs  []uint            `json:"staffTypeIds" openapi:"example:[1,2,3]"`
	StaffTypes    []*CycleStaffType `json:"staffTypes" openapi:"$ref:CycleStaffType;type:array"`
	ShiftName     string            `json:"shiftName" openapi:"example:morning"`
	VehicleType   *string           `json:"vehicleType" openapi:"example:own"`
	StartLocation *string           `json:"startLocation" openapi:"example:office"`
	DateTime      time.Time         `json:"dateTime" openapi:"example:2021-08-02"`
	Status        string            `json:"status" openapi:"example:not-started"`
	CreatedAt     time.Time         `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time         `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time        `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}
