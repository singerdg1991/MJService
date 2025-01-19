package models

/*
 * @apiDefine: StaffsDeleteResponseData
 */
type StaffsDeleteResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: StaffsDeleteResponse
 */
type StaffsDeleteResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200"`
	Data       StaffsDeleteResponseData `json:"data" openapi:"$ref:StaffsDeleteResponseData"`
}
