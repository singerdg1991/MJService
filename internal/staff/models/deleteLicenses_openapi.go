package models

/*
 * @apiDefine: StaffsDeleteLicensesResponseData
 */
type StaffsDeleteLicensesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: StaffsDeleteLicensesResponse
 */
type StaffsDeleteLicensesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       StaffsDeleteLicensesResponseData `json:"data" openapi:"$ref:StaffsDeleteLicensesResponseData"`
}
