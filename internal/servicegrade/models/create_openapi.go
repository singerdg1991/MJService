package models

/*
 * @apiDefine: ServiceGradesResponseData
 */
type ServiceGradesResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
	Grade       int     `json:"grade" openapi:"example:0"`
	Color       string  `json:"color" openapi:"example:#000000"`
}

/*
 * @apiDefine: ServiceGradesCreateResponse
 */
type ServiceGradesCreateResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200"`
	Data       ServiceGradesResponseData `json:"data" openapi:"$ref:ServiceGradesResponseData"`
}
