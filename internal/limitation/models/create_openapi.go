package models

/*
 * @apiDefine: LimitationsResponseData
 */
type LimitationsResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: LimitationsCreateResponse
 */
type LimitationsCreateResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       LimitationsResponseData `json:"data" openapi:"$ref:LimitationsResponseData"`
}
