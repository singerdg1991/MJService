package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateShiftCustomerHomeKeyReleaseResponseData
 */
type CyclesCreateShiftCustomerHomeKeyReleaseResponseData struct {
	ID        uint                                       `json:"id" openapi:"example:1"`
	ShiftID   uint                                       `json:"shiftId" openapi:"example:1"`
	KeyNo     string                                     `json:"keyNo" openapi:"example:1"`
	Status    string                                     `json:"status" openapi:"example:accepted"`
	Reason    *string                                    `json:"reason" openapi:"example:accepted"`
	CreatedAt string                                     `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy *domain.CycleShiftCustomerHomeKeyCreatedBy `json:"createdBy" openapi:"$ref:CycleShiftCustomerHomeKeyCreatedBy;required"`
	UpdatedAt string                                     `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *string                                    `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreateShiftCustomerHomeKeyReleaseResponse
 */
type CyclesCreateShiftCustomerHomeKeyReleaseResponse struct {
	StatusCode int                                                 `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateShiftCustomerHomeKeyReleaseResponseData `json:"data" openapi:"$ref:CyclesCreateShiftCustomerHomeKeyReleaseResponseData"`
}
