package models

/*
 * @apiDefine: AIsCreateResponseData
 */
type AIsCreateResponseData struct {
	Response string `json:"response" openapi:"example:response;required"`
}

/*
 * @apiDefine: AIsCreateResponse
 */
type AIsCreateResponse struct {
	StatusCode int                   `json:"statusCode" openapi:"example:200"`
	Data       AIsCreateResponseData `json:"data" openapi:"$ref:AIsCreateResponseData"`
}
