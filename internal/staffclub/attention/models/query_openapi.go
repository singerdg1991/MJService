package models

import "github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"

/*
 * @apiDefine: AttentionsQueryResponseData
 */
type AttentionsQueryResponseData struct {
	Limit      int                      `json:"limit" openapi:"example:10"`
	Offset     int                      `json:"offset" openapi:"example:0"`
	Page       int                      `json:"page" openapi:"example:1"`
	TotalRows  int                      `json:"totalRows" openapi:"example:1"`
	TotalPages int                      `json:"totalPages" openapi:"example:1"`
	Items      []AttentionsResponseData `json:"items" openapi:"$ref:AttentionsResponseData;type:array"`
}

/*
 * @apiDefine: AttentionsQueryResponse
 */
type AttentionsQueryResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200;"`
	Data       AttentionsQueryResponseData `json:"data" openapi:"$ref:AttentionsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: AttentionsQueryNotFoundResponse
 */
type AttentionsQueryNotFoundResponse struct {
	Attentions []domain.Attention `json:"attentions" openapi:"$ref:Attention;type:array"`
}
