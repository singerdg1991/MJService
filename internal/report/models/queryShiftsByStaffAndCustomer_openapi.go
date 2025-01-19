package models

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryShiftsByStaffAndCustomerChartResponseDataItem
 */
type ReportsQueryShiftsByStaffAndCustomerChartResponseDataItem struct {
	StaffId      int64  `json:"staffId" openapi:"example:1"`
	StaffName    string `json:"staffName" openapi:"example:John Doe"`
	CustomerId   int64  `json:"customerId" openapi:"example:1"`
	CustomerName string `json:"customerName" openapi:"example:Customer A"`
	ShiftCount   int    `json:"shiftCount" openapi:"example:30"`
}

/*
 * @apiDefine: ReportsQueryShiftsByStaffAndCustomerChartResponseData
 */
type ReportsQueryShiftsByStaffAndCustomerChartResponseData struct {
	Items []ReportsQueryShiftsByStaffAndCustomerChartResponseDataItem `json:"items" openapi:"$ref:ReportsQueryShiftsByStaffAndCustomerChartResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryShiftsByStaffAndCustomerChartResponse
 */
type ReportsQueryShiftsByStaffAndCustomerChartResponse struct {
	StatusCode int                                                   `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryShiftsByStaffAndCustomerChartResponseData `json:"data" openapi:"$ref:ReportsQueryShiftsByStaffAndCustomerChartResponseData"`
}

/*
 * @apiDefine: ReportsQueryShiftsByStaffAndCustomerChartNotFoundResponse
 */
type ReportsQueryShiftsByStaffAndCustomerChartNotFoundResponse struct {
	ReportsQueryShiftDurationAnalysis []domain.ReportArrangementTable `json:"ReportsQueryShiftDurationAnalysis" openapi:"$ref:ReportArrangementTable;type:array"`
}
