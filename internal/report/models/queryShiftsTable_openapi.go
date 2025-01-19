package models

import (
	cycleDomain "github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryShiftsTableResponseDataItem
 */
type ReportsQueryShiftsTableResponseDataItem struct {
	ID            uint                          `json:"id" openapi:"example:1"`
	ExchangeKey   string                        `json:"exchangeKey" openapi:"example:dfhdsjrtwerwrwfgjgfrt"`
	CycleID       uint                          `json:"cycleId" openapi:"example:1"`
	StaffTypeIDs  []uint                        `json:"staffTypeIds" openapi:"example:[1,2,3]"`
	StaffTypes    []*cycleDomain.CycleStaffType `json:"staffTypes" openapi:"$ref:CycleStaffType;type:array"`
	ShiftName     string                        `json:"shiftName" openapi:"example:morning"`
	VehicleType   *string                       `json:"vehicleType" openapi:"example:own"`
	StartLocation *string                       `json:"startLocation" openapi:"example:office"`
	DateTime      string                        `json:"dateTime" openapi:"example:2021-08-02"`
	Status        string                        `json:"status" openapi:"example:not-started"`
	CreatedAt     string                        `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     string                        `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *string                       `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: ReportsQueryShiftsTableResponseData
 */
type ReportsQueryShiftsTableResponseData struct {
	Limit      int                                       `json:"limit" openapi:"example:10"`
	Offset     int                                       `json:"offset" openapi:"example:0"`
	Page       int                                       `json:"page" openapi:"example:1"`
	TotalRows  int                                       `json:"totalRows" openapi:"example:1"`
	TotalPages int                                       `json:"totalPages" openapi:"example:1"`
	Items      []ReportsQueryShiftsTableResponseDataItem `json:"items" openapi:"$ref:ReportsQueryShiftsTableResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryShiftsTableResponse
 */
type ReportsQueryShiftsTableResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryShiftsTableResponseData `json:"data" openapi:"$ref:ReportsQueryShiftsTableResponseData"`
}

/*
 * @apiDefine: ReportsQueryShiftsTableNotFoundResponse
 */
type ReportsQueryShiftsTableNotFoundResponse struct {
	ReportsQueryShiftsTable []domain.ReportShiftTable `json:"reportsQueryShiftsTable" openapi:"$ref:ReportShiftTable;type:array"`
}
