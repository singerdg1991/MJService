package models

/*
 * @apiDefine: CustomersDeleteServicesResponseData
 */
type CustomersDeleteServicesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteServicesResponse
 */
type CustomersDeleteServicesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteServicesResponseData `json:"data" openapi:"$ref:CustomersDeleteServicesResponseData"`
}
