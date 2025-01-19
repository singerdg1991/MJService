package models

/*
 * @apiDefine: DiagnosesCreateResponseData
 */
type DiagnosesCreateResponseData struct {
	ID          int    `json:"id" openapi:"example:1"`
	Title       string `json:"title" openapi:"example:title"`
	Code        string `json:"code" openapi:"example:code"`
	Description string `json:"description" openapi:"example:description"`
}

/*
 * @apiDefine: DiagnosesCreateResponse
 */
type DiagnosesCreateResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200;"`
	Data       DiagnosesCreateResponseData `json:"data" openapi:"$ref:DiagnosesCreateResponseData"`
}
