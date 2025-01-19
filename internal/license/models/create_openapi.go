package models

/*
 * @apiDefine: LicensesResponseData
 */
type LicensesResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: LicensesCreateResponse
 */
type LicensesCreateResponse struct {
	StatusCode int                  `json:"statusCode" openapi:"example:200"`
	Data       LicensesResponseData `json:"data" openapi:"$ref:LicensesResponseData"`
}
