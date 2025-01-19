package models

/*
 * @apiDefine: CustomersDeleteCreditDetailsResponseData
 */
type CustomersDeleteCreditDetailsResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: CustomersDeleteCreditDetailsResponse
 */
type CustomersDeleteCreditDetailsResponse struct {
	StatusCode int                                      `json:"statusCode" openapi:"example:200"`
	Data       CustomersDeleteCreditDetailsResponseData `json:"data" openapi:"$ref:CustomersDeleteCreditDetailsResponseData"`
}
