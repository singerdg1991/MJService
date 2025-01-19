package models

import "github.com/hoitek/Maja-Service/internal/servicetype/domain"

/*
 * @apiDefine: ServiceTypesQueryResponseDataItem
 */
type ServiceTypesQueryResponseDataItem struct {
	ID          uint                      `json:"id" openapi:"example:1"`
	ServiceID   uint                      `json:"serviceId" openapi:"example:1"`
	Service     domain.ServiceTypeService `json:"service" openapi:"$ref:ServiceTypeService"`
	Name        string                    `json:"name" openapi:"example:John;required"`
	Description string                    `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServiceTypesQueryResponseData
 */
type ServiceTypesQueryResponseData struct {
	Limit      int                                 `json:"limit" openapi:"example:10"`
	Offset     int                                 `json:"offset" openapi:"example:0"`
	Page       int                                 `json:"page" openapi:"example:1"`
	TotalRows  int                                 `json:"totalRows" openapi:"example:1"`
	TotalPages int                                 `json:"totalPages" openapi:"example:1"`
	Items      []ServiceTypesQueryResponseDataItem `json:"items" openapi:"$ref:ServiceTypesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ServiceTypesQueryResponse
 */
type ServiceTypesQueryResponse struct {
	StatusCode int                           `json:"statusCode" openapi:"example:200;"`
	Data       ServiceTypesQueryResponseData `json:"data" openapi:"$ref:ServiceTypesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ServiceTypesQueryNotFoundResponse
 */
type ServiceTypesQueryNotFoundResponse struct {
	ServiceTypes []domain.ServiceType `json:"serviceTypes" openapi:"$ref:ServiceType;type:array"`
}
