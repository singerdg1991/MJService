package models

/*
 * @apiDefine: StaffsDeleteLibrariesResponseData
 */
type StaffsDeleteLibrariesResponseData struct {
	IDs int `json:"ids" openapi:"example:[1,2];type:array"`
}

/*
 * @apiDefine: StaffsDeleteLibrariesResponse
 */
type StaffsDeleteLibrariesResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       StaffsDeleteLibrariesResponseData `json:"data" openapi:"$ref:StaffsDeleteLibrariesResponseData"`
}
