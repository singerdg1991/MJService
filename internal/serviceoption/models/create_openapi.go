package models

import "github.com/hoitek/Maja-Service/internal/serviceoption/domain"

/*
 * @apiDefine: ServiceOptionsResponseData
 */
type ServiceOptionsResponseData struct {
	ID            uint                            `json:"id" openapi:"example:1"`
	ServiceTypeID uint                            `json:"serviceTypeId" openapi:"example:1"`
	ServiceType   domain.ServiceOptionServiceType `json:"serviceType" openapi:"$ref:ServiceOptionServiceType"`
	Name          string                          `json:"name" openapi:"example:saeed"`
	Description   string                          `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: ServiceOptionsCreateResponse
 */
type ServiceOptionsCreateResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       ServiceOptionsResponseData `json:"data" openapi:"$ref:ServiceOptionsResponseData"`
}
