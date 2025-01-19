package models

import "github.com/hoitek/Maja-Service/internal/permission/domain"

/*
 * @apiDefine: PermissionsQueryResponseDataItem
 */
type PermissionsQueryResponseDataItem struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:John;required"`
	Title string `json:"title" openapi:"example:John;required"`
}

/*
 * @apiDefine: PermissionsQueryResponseData
 */
type PermissionsQueryResponseData struct {
	Limit      int                                `json:"limit" openapi:"example:10"`
	Offset     int                                `json:"offset" openapi:"example:0"`
	Page       int                                `json:"page" openapi:"example:1"`
	TotalRows  int                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                `json:"totalPages" openapi:"example:1"`
	Items      []PermissionsQueryResponseDataItem `json:"items" openapi:"$ref:PermissionsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: PermissionsQueryResponse
 */
type PermissionsQueryResponse struct {
	StatusCode int                          `json:"statusCode" openapi:"example:200;"`
	Data       PermissionsQueryResponseData `json:"data" openapi:"$ref:PermissionsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: PermissionsQueryNotFoundResponse
 */
type PermissionsQueryNotFoundResponse struct {
	Permissions []domain.Permission `json:"permissions" openapi:"$ref:Permission;type:array"`
}
