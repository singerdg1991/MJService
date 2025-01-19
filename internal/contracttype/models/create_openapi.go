package models

/*
 * @apiDefine: ContractTypesResponseData
 */
type ContractTypesResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: ContractTypesCreateResponse
 */
type ContractTypesCreateResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200"`
	Data       ContractTypesResponseData `json:"data" openapi:"$ref:ContractTypesResponseData"`
}
