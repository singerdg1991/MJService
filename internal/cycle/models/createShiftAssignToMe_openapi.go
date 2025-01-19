package models

import (
	csDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateShiftAssignToMeResponseData
 */
type CyclesCreateShiftAssignToMeResponseData struct {
	ID               uint                                   `json:"id" openapi:"example:1"`
	CycleID          uint                                   `json:"cycleId" openapi:"example:1"`
	Staff            *domain.CyclePickupShiftStaff          `json:"staff" openapi:"$ref:CyclePickupShiftStaff"`
	Shift            *domain.CycleShift                     `json:"shift" openapi:"$ref:CycleShift"`
	CycleStaffType   *domain.CyclePickupShiftCycleStaffType `json:"cycleStaffType" openapi:"$ref:CyclePickupShiftCycleStaffType"`
	CustomerServices []*csDomain.CustomerServices           `json:"customerServices" openapi:"$ref:CustomerServices"`
	DateTime         string                                 `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt        string                                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt        string                                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt        *string                                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreateShiftAssignToMeResponse
 */
type CyclesCreateShiftAssignToMeResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateShiftAssignToMeResponseData `json:"data" openapi:"$ref:CyclesCreateShiftAssignToMeResponseData"`
}
