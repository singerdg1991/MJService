package models

/*
 * @apiDefine: PushesSendResponseData
 */
type PushesSendResponseData struct {
	Status string `json:"status" openapi:"example:success"`
}

/*
 * @apiDefine: PushesSendResponse
 */
type PushesSendResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200"`
	Data       PushesSendResponseData `json:"data" openapi:"$ref:PushesSendResponseData"`
}
