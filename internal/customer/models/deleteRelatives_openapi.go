package models

/*
 * @apiDefine: CustomersDeleteRelativesResponseData
 */
type CustomersDeleteRelativesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteRelativesResponse
 */
type CustomersDeleteRelativesResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteRelativesResponseData `json:"data" openapi:"$ref:CustomersDeleteRelativesResponseData"`
}
