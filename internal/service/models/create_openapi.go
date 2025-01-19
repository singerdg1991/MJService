package models

/*
 * @apiDefine: ServicesResponseData
 */
type ServicesResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: ServicesCreateResponse
 */
type ServicesCreateResponse struct {
	StatusCode int                  `json:"statusCode" openapi:"example:200"`
	Data       ServicesResponseData `json:"data" openapi:"$ref:ServicesResponseData"`
}
