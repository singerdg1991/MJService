package models

/*
 * @apiDefine: StaffsDeleteAbsencesResponseData
 */
type StaffsDeleteAbsencesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: StaffsDeleteAbsencesResponse
 */
type StaffsDeleteAbsencesResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       StaffsDeleteAbsencesResponseData `json:"data" openapi:"$ref:StaffsDeleteAbsencesResponseData"`
}
