package models

import "github.com/hoitek/Maja-Service/internal/trash/domain"

/*
 * @apiDefine: TrashesQueryResponseData
 */
type TrashesQueryResponseData struct {
	Limit      int            `json:"limit" openapi:"example:10"`
	Offset     int            `json:"offset" openapi:"example:0"`
	Page       int            `json:"page" openapi:"example:1"`
	TotalRows  int            `json:"totalRows" openapi:"example:1"`
	TotalPages int            `json:"totalPages" openapi:"example:1"`
	Items      []domain.Trash `json:"trashes" openapi:"$ref:Trash;type:array"`
}

/*
 * @apiDefine: TrashesQueryResponse
 */
type TrashesQueryResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200"`
	Data       TrashesQueryResponseData `json:"data" openapi:"$ref:TrashesQueryResponseData"`
}

/*
 * @apiDefine: TrashesQueryNotFoundResponse
 */
type TrashesQueryNotFoundResponse struct {
	Trashes []domain.Trash `json:"trashes" openapi:"$ref:Trash;type:array"`
}
