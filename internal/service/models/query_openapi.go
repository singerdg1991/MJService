package models

import "github.com/hoitek/Maja-Service/internal/service/domain"

/*
 * @apiDefine: ServicesQueryResponseDataItem
 */
type ServicesQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServicesQueryResponseData
 */
type ServicesQueryResponseData struct {
	Limit      int                             `json:"limit" openapi:"example:10"`
	Offset     int                             `json:"offset" openapi:"example:0"`
	Page       int                             `json:"page" openapi:"example:1"`
	TotalRows  int                             `json:"totalRows" openapi:"example:1"`
	TotalPages int                             `json:"totalPages" openapi:"example:1"`
	Items      []ServicesQueryResponseDataItem `json:"items" openapi:"$ref:ServicesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ServicesQueryResponse
 */
type ServicesQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       ServicesQueryResponseData `json:"data" openapi:"$ref:ServicesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ServicesQueryNotFoundResponse
 */
type ServicesQueryNotFoundResponse struct {
	Services []domain.Service `json:"services" openapi:"$ref:Service;type:array"`
}
