package models

import "github.com/hoitek/Maja-Service/internal/stafftype/domain"

/*
 * @apiDefine: StaffTypesQueryResponseDataItem
 */
type StaffTypesQueryResponseDataItem struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:name"`
}

/*
 * @apiDefine: StaffTypesQueryResponseData
 */
type StaffTypesQueryResponseData struct {
	Limit      int                               `json:"limit" openapi:"example:10"`
	Offset     int                               `json:"offset" openapi:"example:0"`
	Page       int                               `json:"page" openapi:"example:1"`
	TotalRows  int                               `json:"totalRows" openapi:"example:1"`
	TotalPages int                               `json:"totalPages" openapi:"example:1"`
	Items      []StaffTypesQueryResponseDataItem `json:"items" openapi:"$ref:StaffTypesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: StaffTypesQueryResponse
 */
type StaffTypesQueryResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200"`
	Data       StaffTypesQueryResponseData `json:"data" openapi:"$ref:StaffTypesQueryResponseData"`
}

/*
 * @apiDefine: StaffTypesQueryNotFoundResponse
 */
type StaffTypesQueryNotFoundResponse struct {
	StaffTypes []domain.StaffType `json:"staffTypes" openapi:"$ref:StaffType;type:array"`
}
