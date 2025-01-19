package models

/*
 * @apiDefine: RewardsResponseData
 */
type RewardsResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: RewardsCreateResponse
 */
type RewardsCreateResponse struct {
	StatusCode int                 `json:"statusCode" openapi:"example:200"`
	Data       RewardsResponseData `json:"data" openapi:"$ref:RewardsResponseData"`
}
