package models

import "github.com/hoitek/Maja-Service/internal/report/domain"

/*
 * @apiDefine: ReportsQueryShiftHeatmapChartResponseData
 */
type ReportsQueryShiftHeatmapChartResponseData struct {
	DayOfWeek  string `json:"dayOfWeek" openapi:"example:Monday"`
	DayOrder   int    `json:"dayOrder" openapi:"example:1"` // 1 for Monday
	Hour       int    `json:"hour" openapi:"example:9"`     // 9 AM
	ShiftCount int    `json:"shiftCount" openapi:"example:5"`
}

/*
 * @apiDefine: ReportsQueryShiftHeatmapChartResponse
 */
type ReportsQueryShiftHeatmapChartResponse struct {
	StatusCode int                                         `json:"statusCode" openapi:"example:200"`
	Data       []ReportsQueryShiftHeatmapChartResponseData `json:"data" openapi:"$ref:ReportsQueryShiftHeatmapChartResponseData;type:array"`
}

/*
 * @apiDefine: ReportsQueryShiftHeatmapChartNotFoundResponse
 */
type ReportsQueryShiftHeatmapChartNotFoundResponse struct {
	ReportsQueryShiftHeatmap []domain.ReportArrangementTable `json:"ReportsQueryShiftHeatmap" openapi:"$ref:ReportArrangementTable;type:array"`
}
