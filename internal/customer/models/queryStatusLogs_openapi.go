package models

import (
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersQueryStatusLogsResponseDataItem
 */
type CustomersQueryStatusLogsResponseDataItem struct {
	ID          int64                              `json:"id" openapi:"example:1"`
	CustomerID  int64                              `json:"customerId" openapi:"example:1"`
	StatusValue string                             `json:"statusValue" openapi:"example:StatusValue"`
	Customer    *domain.CustomerStatusLogCustomer  `json:"customer" openapi:"$ref:CustomerStatusLogCustomer"`
	CreatedBy   *domain.CustomerStatusLogCreatedBy `json:"createdBy" openapi:"$ref:CustomerStatusLogCreatedBy"`
	CreatedAt   string                             `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   string                             `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *string                            `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersQueryStatusLogsResponseData
 */
type CustomersQueryStatusLogsResponseData struct {
	Limit      int                                        `json:"limit" openapi:"example:10"`
	Offset     int                                        `json:"offset" openapi:"example:0"`
	Page       int                                        `json:"page" openapi:"example:1"`
	TotalRows  int                                        `json:"totalRows" openapi:"example:1"`
	TotalPages int                                        `json:"totalPages" openapi:"example:1"`
	Items      []CustomersQueryStatusLogsResponseDataItem `json:"items" openapi:"$ref:CustomersQueryStatusLogsResponseDataItem"`
}

/*
 * @apiDefine: CustomersQueryStatusLogsResponse
 */
type CustomersQueryStatusLogsResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryStatusLogsResponseData `json:"data" openapi:"$ref:CustomersQueryStatusLogsResponseData"`
}
