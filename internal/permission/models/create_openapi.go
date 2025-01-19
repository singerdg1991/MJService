package models

/*
 * @apiDefine: PermissionsResponseData
 */
type PermissionsResponseData struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:saeed"`
	Title string `json:"title" openapi:"example:test"`
}

/*
 * @apiDefine: PermissionsCreateResponse
 */
type PermissionsCreateResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       PermissionsResponseData `json:"data" openapi:"$ref:PermissionsResponseData"`
}
