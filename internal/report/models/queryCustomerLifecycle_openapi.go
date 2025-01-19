package models

/*
 * @apiDefine: ReportsQueryCustomerLifecycleChartResponseData
 */
type ReportsQueryCustomerLifecycleChartResponseData struct {
	Stage      string `json:"stage" openapi:"example:active"`
	Count      int    `json:"count" openapi:"example:100"`
	StageOrder int    `json:"stageOrder" openapi:"example:1"`
}

/*
 * @apiDefine: ReportsQueryCustomerLifecycleChartResponse
 */
type ReportsQueryCustomerLifecycleChartResponse struct {
	StatusCode int                                              `json:"statusCode" openapi:"example:200"`
	Data       []ReportsQueryCustomerLifecycleChartResponseData `json:"data" openapi:"$ref:ReportsQueryCustomerLifecycleChartResponseData;type:array"`
}

/*
 * @apiDefine: ReportsQueryCustomerLifecycleChartNotFoundResponse
 */
type ReportsQueryCustomerLifecycleChartNotFoundResponse struct {
	ReportsQueryCustomerLifecycleChartResponseData []ReportsQueryCustomerLifecycleChartResponseData `json:"ReportsQueryCustomerLifecycleChartResponseData" openapi:"$ref:ReportsQueryCustomerLifecycleChartResponseData;type:array"`
}
