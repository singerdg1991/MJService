package models

import "github.com/hoitek/Maja-Service/internal/servicetype/domain"

/*
 * @apiDefine: ServiceTypesResponseData
 */
type ServiceTypesResponseData struct {
	ID          uint                      `json:"id" openapi:"example:1"`
	ServiceID   uint                      `json:"serviceId" openapi:"example:1"`
	Service     domain.ServiceTypeService `json:"service" openapi:"$ref:ServiceTypeService"`
	Name        string                    `json:"name" openapi:"example:saeed"`
	Description string                    `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: ServiceTypesCreateResponse
 */
type ServiceTypesCreateResponse struct {
	StatusCode int                      `json:"statusCode" openapi:"example:200"`
	Data       ServiceTypesResponseData `json:"data" openapi:"$ref:ServiceTypesResponseData"`
}
