package models

import "github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"

/*
 * @apiDefine: WarningsQueryResponseData
 */
type WarningsQueryResponseData struct {
	Limit      int                    `json:"limit" openapi:"example:10"`
	Offset     int                    `json:"offset" openapi:"example:0"`
	Page       int                    `json:"page" openapi:"example:1"`
	TotalRows  int                    `json:"totalRows" openapi:"example:1"`
	TotalPages int                    `json:"totalPages" openapi:"example:1"`
	Items      []WarningsResponseData `json:"items" openapi:"$ref:WarningsResponseData;type:array"`
}

/*
 * @apiDefine: WarningsQueryResponse
 */
type WarningsQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       WarningsQueryResponseData `json:"data" openapi:"$ref:WarningsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: WarningsQueryNotFoundResponse
 */
type WarningsQueryNotFoundResponse struct {
	Warnings []domain.Warning `json:"warnings" openapi:"$ref:Warning;type:array"`
}
