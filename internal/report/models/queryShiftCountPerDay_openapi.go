package models

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryShiftCountPerDayChartResponseData
 */
type ReportsQueryShiftCountPerDayChartResponseData struct {
	DayOfWeek  string `json:"dayOfWeek" openapi:"example:Monday"`
	ShiftCount int    `json:"shiftCount" openapi:"example:25"`
	DayOrder   int    `json:"dayOrder" openapi:"example:1"` // 1 for Monday, 2 for Tuesday, etc.
}

/*
 * @apiDefine: ReportsQueryShiftCountPerDayChartResponse
 */
type ReportsQueryShiftCountPerDayChartResponse struct {
	StatusCode int                                           `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryShiftCountPerDayChartResponseData `json:"data" openapi:"$ref:ReportsQueryShiftCountPerDayChartResponseData"`
}

/*
 * @apiDefine: ReportsQueryShiftCountPerDayChartNotFoundResponse
 */
type ReportsQueryShiftCountPerDayChartNotFoundResponse struct {
	ReportsQueryShiftCountPerDay []domain.ReportArrangementTable `json:"ReportsQueryShiftCountPerDay" openapi:"$ref:ReportArrangementTable;type:array"`
}
