package models

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesUpdateNextStaffTypesResponseData
 */
type CyclesUpdateNextStaffTypesResponseData struct {
	Limit      int                                         `json:"limit" openapi:"example:10"`
	Offset     int                                         `json:"offset" openapi:"example:0"`
	Page       int                                         `json:"page" openapi:"example:1"`
	TotalRows  int                                         `json:"totalRows" openapi:"example:1"`
	TotalPages int                                         `json:"totalPages" openapi:"example:1"`
	Items      []CyclesUpdateNextStaffTypeResponseDataItem `json:"items" openapi:"$ref:CyclesUpdateNextStaffTypeResponseDataItem;type:array"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypesResponse
 */
type CyclesUpdateNextStaffTypesResponse struct {
	StatusCode int                                    `json:"statusCode" openapi:"example:200;"`
	Date       CyclesUpdateNextStaffTypesResponseData `json:"data" openapi:"$ref:CyclesUpdateNextStaffTypesResponseData;type:object;"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypesNotFoundResponse
 */
type CyclesUpdateNextStaffTypesNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
