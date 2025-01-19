package models

import "github.com/hoitek/Maja-Service/internal/push/domain"

/*
 * @apiDefine: PushesQueryResponseData
 */
type PushesQueryResponseData struct {
	Limit      int                  `json:"limit" openapi:"example:10"`
	Offset     int                  `json:"offset" openapi:"example:0"`
	Page       int                  `json:"page" openapi:"example:1"`
	TotalRows  int                  `json:"totalRows" openapi:"example:1"`
	TotalPages int                  `json:"totalPages" openapi:"example:1"`
	Items      []PushesResponseData `json:"items" openapi:"$ref:PushesResponseData;type:array"`
}

/*
 * @apiDefine: PushesQueryResponse
 */
type PushesQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200;"`
	Data       PushesQueryResponseData `json:"data" openapi:"$ref:PushesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: PushesQueryNotFoundResponse
 */
type PushesQueryNotFoundResponse struct {
	Pushes []domain.Push `json:"pushes" openapi:"$ref:Push;type:array"`
}
