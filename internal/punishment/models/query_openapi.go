package models

import "github.com/hoitek/Maja-Service/internal/punishment/domain"

/*
 * @apiDefine: PunishmentsQueryResponseDataItem
 */
type PunishmentsQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: PunishmentsQueryResponseData
 */
type PunishmentsQueryResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []PunishmentsQueryResponseDataItem `json:"items" openapi:"$ref:PunishmentsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: PunishmentsQueryResponse
 */
type PunishmentsQueryResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200;"`
	Data       PunishmentsQueryResponseData `json:"data" openapi:"$ref:PunishmentsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: PunishmentsQueryNotFoundResponse
 */
type PunishmentsQueryNotFoundResponse struct {
	Punishments []domain.Punishment `json:"punishments" openapi:"$ref:Punishment;type:array"`
}
