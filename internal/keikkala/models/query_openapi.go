package models

import "github.com/hoitek/Maja-Service/internal/keikkala/domain"

/*
 * @apiDefine: KeikkalasQueryResponseData
 */
type KeikkalasQueryResponseData struct {
	Limit      int                     `json:"limit" openapi:"example:10"`
	Offset     int                     `json:"offset" openapi:"example:0"`
	Page       int                     `json:"page" openapi:"example:1"`
	TotalRows  int                     `json:"totalRows" openapi:"example:1"`
	TotalPages int                     `json:"totalPages" openapi:"example:1"`
	Items      []KeikkalasResponseData `json:"items" openapi:"$ref:KeikkalasResponseData;type:array"`
}

/*
 * @apiDefine: KeikkalasQueryResponse
 */
type KeikkalasQueryResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200;"`
	Data       KeikkalasQueryResponseData `json:"data" openapi:"$ref:KeikkalasQueryResponseData;type:object;"`
}

/*
 * @apiDefine: KeikkalasQueryNotFoundResponse
 */
type KeikkalasQueryNotFoundResponse struct {
	Keikkalas []domain.Keikkala `json:"keikkalas" openapi:"$ref:Keikkala;type:array"`
}
