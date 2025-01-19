package models

/*
 * @apiDefine: EquipmentsResponseData
 */
type EquipmentsResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: EquipmentsCreateResponse
 */
type EquipmentsCreateResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200"`
	Data       EquipmentsResponseData `json:"data" openapi:"$ref:EquipmentsResponseData"`
}
