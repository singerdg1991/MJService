package models

import "github.com/hoitek/Maja-Service/internal/limitation/domain"

/*
 * @apiDefine: LimitationsQueryResponseDataItem
 */
type LimitationsQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: LimitationsQueryResponseData
 */
type LimitationsQueryResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []LimitationsQueryResponseDataItem `json:"items" openapi:"$ref:LimitationsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: LimitationsQueryResponse
 */
type LimitationsQueryResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200;"`
	Data       LimitationsQueryResponseData `json:"data" openapi:"$ref:LimitationsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: LimitationsQueryNotFoundResponse
 */
type LimitationsQueryNotFoundResponse struct {
	Limitations []domain.Limitation `json:"limitations" openapi:"$ref:Limitation;type:array"`
}
