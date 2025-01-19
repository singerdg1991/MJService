package models

import "github.com/hoitek/Maja-Service/internal/role/domain"

/*
 * @apiDefine: RolesQueryResponseData
 */
type RolesQueryResponseData struct {
	Limit      int                       `json:"limit" openapi:"example:10"`
	Offset     int                       `json:"offset" openapi:"example:0"`
	Page       int                       `json:"page" openapi:"example:1"`
	TotalRows  int                       `json:"totalRows" openapi:"example:1"`
	TotalPages int                       `json:"totalPages" openapi:"example:1"`
	Items      []RolesCreateResponseData `json:"items" openapi:"$ref:RolesCreateResponseData;type:array"`
}

/*
 * @apiDefine: RolesQueryResponse
 */
type RolesQueryResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200"`
	Data       RolesQueryResponseData `json:"data" openapi:"$ref:RolesQueryResponseData"`
}

/*
 * @apiDefine: RolesQueryNotFoundResponse
 */
type RolesQueryNotFoundResponse struct {
	Roles []domain.Role `json:"roles" openapi:"$ref:Role;type:array"`
}
