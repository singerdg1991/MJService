package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateShiftCustomerHomeKeyResponseData
 */
type CyclesCreateShiftCustomerHomeKeyResponseData struct {
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
 * @apiDefine: CyclesCreateShiftCustomerHomeKeyResponse
 */
type CyclesCreateShiftCustomerHomeKeyResponse struct {
	StatusCode int                                          `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateShiftCustomerHomeKeyResponseData `json:"data" openapi:"$ref:CyclesCreateShiftCustomerHomeKeyResponseData"`
}
