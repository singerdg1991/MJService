package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"time"
)

/*
 * @apiDefine: CyclesQueryShiftsCustomersResponseDataItem
 */
type CyclesQueryShiftsCustomersResponseDataItem struct {
	ID          uint                               `json:"id" openapi:"example:1"`
	CycleID     uint                               `json:"cycleId" openapi:"example:1"`
	StaffTypeID uint                               `json:"staffTypeId" openapi:"example:1"`
	CustomerID  uint                               `json:"customerId" openapi:"example:1"`
	Customer    *domain.CycleShiftCustomerCustomer `json:"customer" openapi:"$ref:CycleShiftCustomerCustomer"`
	CreatedAt   time.Time                          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time                         `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesQueryShiftsCustomersResponseData
 */
type CyclesQueryShiftsCustomersResponseData struct {
	Limit      int                                          `json:"limit" openapi:"example:10"`
	Offset     int                                          `json:"offset" openapi:"example:0"`
	Page       int                                          `json:"page" openapi:"example:1"`
	TotalRows  int                                          `json:"totalRows" openapi:"example:1"`
	TotalPages int                                          `json:"totalPages" openapi:"example:1"`
	Items      []CyclesQueryShiftsCustomersResponseDataItem `json:"items" openapi:"$ref:CyclesQueryShiftsCustomersResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesQueryShiftsCustomersResponse
 */
type CyclesQueryShiftsCustomersResponse struct {
	StatusCode int                                    `json:"statusCode" openapi:"example:200;"`
	Data       CyclesQueryShiftsCustomersResponseData `json:"data" openapi:"$ref:CyclesQueryShiftsCustomersResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesQueryShiftsCustomersNotFoundResponse
 */
type CyclesQueryShiftsCustomersNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
