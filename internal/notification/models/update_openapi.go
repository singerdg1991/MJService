package models

/*
 * @apiDefine: NotificationsUpdateStatusResponse
 */
type NotificationsUpdateStatusResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       NotificationsResponseData `json:"data" openapi:"$ref:NotificationsResponseData;type:object;"`
}
