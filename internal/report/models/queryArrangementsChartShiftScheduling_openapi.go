package models

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryArrangementsChartShiftSchedulingResponseDataItem
 */
type ReportsQueryArrangementsChartShiftSchedulingResponseDataItem struct {
	StaffID     uint   `json:"staffId" openapi:"example:1"`
	StaffName   string `json:"staffName" openapi:"example:John Doe"`
	ShiftName   string `json:"shiftName" openapi:"example:morning"`
	StartHour   string `json:"startHour" openapi:"example:2021-01-01T09:00:00Z"`
	EndHour     string `json:"endHour" openapi:"example:2021-01-01T17:00:00Z"`
	Date        string `json:"date" openapi:"example:2021-01-01T00:00:00Z"`
	Status      string `json:"status" openapi:"example:assigned"`
	IsUnplanned bool   `json:"isUnplanned" openapi:"example:false"`
}

/*
 * @apiDefine: ReportsQueryArrangementsChartShiftSchedulingResponseData
 */
type ReportsQueryArrangementsChartShiftSchedulingResponseData struct {
	Limit      int                                                            `json:"limit" openapi:"example:10"`
	Offset     int                                                            `json:"offset" openapi:"example:0"`
	Page       int                                                            `json:"page" openapi:"example:1"`
	TotalRows  int                                                            `json:"totalRows" openapi:"example:1"`
	TotalPages int                                                            `json:"totalPages" openapi:"example:1"`
	Items      []ReportsQueryArrangementsChartShiftSchedulingResponseDataItem `json:"items" openapi:"$ref:ReportsQueryArrangementsChartShiftSchedulingResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryArrangementsChartShiftSchedulingResponse
 */
type ReportsQueryArrangementsChartShiftSchedulingResponse struct {
	StatusCode int                                                      `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryArrangementsChartShiftSchedulingResponseData `json:"data" openapi:"$ref:ReportsQueryArrangementsChartShiftSchedulingResponseData"`
}

/*
 * @apiDefine: ReportsQueryArrangementsChartShiftSchedulingNotFoundResponse
 */
type ReportsQueryArrangementsChartShiftSchedulingNotFoundResponse struct {
	ReportsQueryArrangementsChartShiftScheduling []domain.ReportArrangementTable `json:"reportsQueryArrangementsChartShiftScheduling" openapi:"$ref:ReportArrangementTable;type:array"`
}
