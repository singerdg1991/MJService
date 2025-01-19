package models

import "github.com/hoitek/Maja-Service/internal/equipment/domain"

/*
 * @apiDefine: EquipmentsQueryResponseDataItem
 */
type EquipmentsQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: EquipmentsQueryResponseData
 */
type EquipmentsQueryResponseData struct {
	Limit      int                               `json:"limit" openapi:"example:10"`
	Offset     int                               `json:"offset" openapi:"example:0"`
	Page       int                               `json:"page" openapi:"example:1"`
	TotalRows  int                               `json:"totalRows" openapi:"example:1"`
	TotalPages int                               `json:"totalPages" openapi:"example:1"`
	Items      []EquipmentsQueryResponseDataItem `json:"items" openapi:"$ref:EquipmentsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: EquipmentsQueryResponse
 */
type EquipmentsQueryResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200;"`
	Data       EquipmentsQueryResponseData `json:"data" openapi:"$ref:EquipmentsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: EquipmentsQueryNotFoundResponse
 */
type EquipmentsQueryNotFoundResponse struct {
	Equipments []domain.Equipment `json:"equipments" openapi:"$ref:Equipment;type:array"`
}
