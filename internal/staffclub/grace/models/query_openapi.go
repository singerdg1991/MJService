package models

import "github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"

/*
 * @apiDefine: GracesQueryResponseData
 */
type GracesQueryResponseData struct {
	Limit      int                  `json:"limit" openapi:"example:10"`
	Offset     int                  `json:"offset" openapi:"example:0"`
	Page       int                  `json:"page" openapi:"example:1"`
	TotalRows  int                  `json:"totalRows" openapi:"example:1"`
	TotalPages int                  `json:"totalPages" openapi:"example:1"`
	Items      []GracesResponseData `json:"items" openapi:"$ref:GracesResponseData;type:array"`
}

/*
 * @apiDefine: GracesQueryResponse
 */
type GracesQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200;"`
	Data       GracesQueryResponseData `json:"data" openapi:"$ref:GracesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: GracesQueryNotFoundResponse
 */
type GracesQueryNotFoundResponse struct {
	Graces []domain.Grace `json:"graces" openapi:"$ref:Grace;type:array"`
}
