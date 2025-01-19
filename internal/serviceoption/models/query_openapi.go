package models

import "github.com/hoitek/Maja-Service/internal/serviceoption/domain"

/*
 * @apiDefine: ServiceOptionsQueryResponseDataItem
 */
type ServiceOptionsQueryResponseDataItem struct {
	ID            uint                            `json:"id" openapi:"example:1"`
	ServiceTypeID uint                            `json:"serviceTypeId" openapi:"example:1"`
	ServiceType   domain.ServiceOptionServiceType `json:"serviceType" openapi:"$ref:ServiceOptionServiceType"`
	Name          string                          `json:"name" openapi:"example:John;required"`
	Description   string                          `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: ServiceOptionsQueryResponseData
 */
type ServiceOptionsQueryResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []ServiceOptionsQueryResponseDataItem `json:"items" openapi:"$ref:ServiceOptionsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ServiceOptionsQueryResponse
 */
type ServiceOptionsQueryResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200;"`
	Data       ServiceOptionsQueryResponseData `json:"data" openapi:"$ref:ServiceOptionsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ServiceOptionsQueryNotFoundResponse
 */
type ServiceOptionsQueryNotFoundResponse struct {
	ServiceOptions []domain.ServiceOption `json:"serviceOptions" openapi:"$ref:ServiceOption;type:array"`
}
