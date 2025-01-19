package models

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryShiftDistributionByCustomerChartDistributionResponseDataItem
 */
type ReportsQueryShiftDistributionByCustomerChartDistributionResponseDataItem struct {
	CustomerID     uint    `json:"customerId" openapi:"example:1"`
	CustomerName   string  `json:"customerName" openapi:"example:Customer A"`
	NumberOfShifts int     `json:"numberOfShifts" openapi:"example:30"`
	Percentage     float64 `json:"percentage" openapi:"example:25.5"`
}

/*
 * @apiDefine: ReportsQueryShiftDistributionByCustomerChartDistributionResponseData
 */
type ReportsQueryShiftDistributionByCustomerChartDistributionResponseData struct {
	Items []ReportsQueryShiftDistributionByCustomerChartDistributionResponseDataItem `json:"items" openapi:"$ref:ReportsQueryShiftDistributionByCustomerChartDistributionResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryShiftDistributionByCustomerChartDistributionResponse
 */
type ReportsQueryShiftDistributionByCustomerChartDistributionResponse struct {
	StatusCode int                                                                  `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryShiftDistributionByCustomerChartDistributionResponseData `json:"data" openapi:"$ref:ReportsQueryShiftDistributionByCustomerChartDistributionResponseData"`
}

/*
 * @apiDefine: ReportsQueryShiftDistributionByCustomerChartDistributionNotFoundResponse
 */
type ReportsQueryShiftDistributionByCustomerChartDistributionNotFoundResponse struct {
	ReportsQueryShiftDistributionByCustomerChartDistribution []domain.ReportArrangementTable `json:"ReportsQueryShiftDistributionByCustomerChartDistribution" openapi:"$ref:ReportArrangementTable;type:array"`
}
