package models

/*
 * @apiDefine: ReportsQueryCustomerLocationChartResponseData
 */
type ReportsQueryCustomerLocationChartResponseData struct {
	CityID        int    `json:"cityId" openapi:"example:1"`
	CityName      string `json:"cityName" openapi:"example:Helsinki"`
	CustomerCount int    `json:"customerCount" openapi:"example:150"`
}

/*
 * @apiDefine: ReportsQueryCustomerLocationChartResponse
 */
type ReportsQueryCustomerLocationChartResponse struct {
	StatusCode int                                             `json:"statusCode" openapi:"example:200"`
	Data       []ReportsQueryCustomerLocationChartResponseData `json:"data" openapi:"$ref:ReportsQueryCustomerLocationChartResponseData;type:array"`
}

/*
 * @apiDefine: ReportsQueryCustomerLocationChartNotFoundResponse
 */
type ReportsQueryCustomerLocationChartNotFoundResponse struct {
	ReportsQueryCustomerLocationChartResponseData []ReportsQueryCustomerLocationChartResponseData `json:"ReportsQueryCustomerLocationChartResponseData" openapi:"$ref:ReportsQueryCustomerLocationChartResponseData;type:array"`
}
