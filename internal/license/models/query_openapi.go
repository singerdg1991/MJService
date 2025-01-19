package models

import "github.com/hoitek/Maja-Service/internal/license/domain"

/*
 * @apiDefine: LicensesQueryResponseDataItem
 */
type LicensesQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: LicensesQueryResponseData
 */
type LicensesQueryResponseData struct {
	Limit      int                             `json:"limit" openapi:"example:10"`
	Offset     int                             `json:"offset" openapi:"example:0"`
	Page       int                             `json:"page" openapi:"example:1"`
	TotalRows  int                             `json:"totalRows" openapi:"example:1"`
	TotalPages int                             `json:"totalPages" openapi:"example:1"`
	Items      []LicensesQueryResponseDataItem `json:"items" openapi:"$ref:LicensesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: LicensesQueryResponse
 */
type LicensesQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       LicensesQueryResponseData `json:"data" openapi:"$ref:LicensesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: LicensesQueryNotFoundResponse
 */
type LicensesQueryNotFoundResponse struct {
	Licenses []domain.License `json:"licenses" openapi:"$ref:License;type:array"`
}
