package models

/*
 * @apiDefine: ReportsQueryCustomerAgeGroupChartResponseData
 */
type ReportsQueryCustomerAgeGroupChartResponseData struct {
	AgeGroup      string `json:"ageGroup" openapi:"example:21-40"`
	CustomerCount int    `json:"customerCount" openapi:"example:150"`
	GroupOrder    int    `json:"groupOrder" openapi:"example:2"`
}

/*
 * @apiDefine: ReportsQueryCustomerAgeGroupChartResponse
 */
type ReportsQueryCustomerAgeGroupChartResponse struct {
	StatusCode int                                             `json:"statusCode" openapi:"example:200"`
	Data       []ReportsQueryCustomerAgeGroupChartResponseData `json:"data" openapi:"$ref:ReportsQueryCustomerAgeGroupChartResponseData;type:array"`
}

/*
 * @apiDefine: ReportsQueryCustomerAgeGroupChartNotFoundResponse
 */
type ReportsQueryCustomerAgeGroupChartNotFoundResponse struct {
	ReportsQueryCustomerAgeGroupChartResponseData []ReportsQueryCustomerAgeGroupChartResponseData `json:"ReportsQueryCustomerAgeGroupChartResponseData" openapi:"$ref:ReportsQueryCustomerAgeGroupChartResponseData;type:array"`
}
