package models

/*
 * @apiDefine: CustomersDeleteAbsencesResponseData
 */
type CustomersDeleteAbsencesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteAbsencesResponse
 */
type CustomersDeleteAbsencesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteAbsencesResponseData `json:"data" openapi:"$ref:CustomersDeleteAbsencesResponseData"`
}
