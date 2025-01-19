package models

/*
 * @apiDefine: StaffsUpdateLibraryResponse
 */
type StaffsUpdateLibraryResponse struct {
	StatusCode int                               `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateLibrariesResponseData `json:"data" openapi:"$ref:StaffsCreateLibrariesResponseData"`
}
