package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateIncomingCyclePickupShiftResponseData
 */
type CyclesCreateIncomingCyclePickupShiftResponseData struct {
	ID                 uint                                                    `json:"id" openapi:"example:1"`
	CycleID            uint                                                    `json:"cycleId" openapi:"example:1"`
	Staff              *domain.CycleIncomingCyclePickupShiftStaff              `json:"staff" openapi:"$ref:CycleIncomingCyclePickupShiftStaff"`
	CycleNextStaffType *domain.CycleIncomingCyclePickupShiftCycleNextStaffType `json:"cycleNextStaffType" openapi:"$ref:CycleIncomingCyclePickupShiftCycleNextStaffType"`
	DateTime           string                                                  `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt          string                                                  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt          string                                                  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt          *string                                                 `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreateIncomingCyclePickupShiftResponse
 */
type CyclesCreateIncomingCyclePickupShiftResponse struct {
	StatusCode int                                              `json:"statusCode" openapi:"example:200"`
	Data       CyclesQueryIncomingCyclePickupShiftsResponseData `json:"data" openapi:"$ref:CyclesQueryIncomingCyclePickupShiftsResponseData"`
}
