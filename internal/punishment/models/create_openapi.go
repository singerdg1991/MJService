package models

/*
 * @apiDefine: PunishmentsResponseData
 */
type PunishmentsResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: PunishmentsCreateResponse
 */
type PunishmentsCreateResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       PunishmentsResponseData `json:"data" openapi:"$ref:PunishmentsResponseData"`
}
