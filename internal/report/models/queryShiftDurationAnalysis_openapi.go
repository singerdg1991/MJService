package models

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryShiftDurationAnalysisResponseDataItem
 */
type ReportsQueryShiftDurationAnalysisResponseDataItem struct {
	CustomerID     uint    `json:"customerId" openapi:"example:1"`
	CustomerName   string  `json:"customerName" openapi:"example:Customer A"`
	DurationHours  float64 `json:"durationHours" openapi:"example:6"`
	NumberOfShifts int     `json:"numberOfShifts" openapi:"example:30"`
}

/*
 * @apiDefine: ReportsQueryShiftDurationAnalysisResponseData
 */
type ReportsQueryShiftDurationAnalysisResponseData struct {
	Items []ReportsQueryShiftDurationAnalysisResponseDataItem `json:"items" openapi:"$ref:ReportsQueryShiftDurationAnalysisResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryShiftDurationAnalysisResponse
 */
type ReportsQueryShiftDurationAnalysisResponse struct {
	StatusCode int                                           `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryShiftDurationAnalysisResponseData `json:"data" openapi:"$ref:ReportsQueryShiftDurationAnalysisResponseData"`
}

/*
 * @apiDefine: ReportsQueryShiftDurationAnalysisNotFoundResponse
 */
type ReportsQueryShiftDurationAnalysisNotFoundResponse struct {
	ReportsQueryShiftDurationAnalysis []domain.ReportArrangementTable `json:"ReportsQueryShiftDurationAnalysis" openapi:"$ref:ReportArrangementTable;type:array"`
}
