package models

import (
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsResponseDataItem
 */
type CustomersQueryContractualMobilityRestrictionLogsResponseDataItem struct {
	ID          int64                                                      `json:"id" openapi:"example:1"`
	CustomerID  int64                                                      `json:"customerId" openapi:"example:1"`
	BeforeValue *string                                                    `json:"beforeValue" openapi:"example:BeforeValue"`
	AfterValue  *string                                                    `json:"afterValue" openapi:"example:AfterValue"`
	Customer    *domain.CustomerContractualMobilityRestrictionLogCustomer  `json:"customer" openapi:"$ref:CustomerContractualMobilityRestrictionLogCustomer"`
	CreatedBy   *domain.CustomerContractualMobilityRestrictionLogCreatedBy `json:"createdBy" openapi:"$ref:CustomerContractualMobilityRestrictionLogCreatedBy"`
	CreatedAt   string                                                     `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   string                                                     `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *string                                                    `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsResponseData
 */
type CustomersQueryContractualMobilityRestrictionLogsResponseData struct {
	Limit      int                                                                `json:"limit" openapi:"example:10"`
	Offset     int                                                                `json:"offset" openapi:"example:0"`
	Page       int                                                                `json:"page" openapi:"example:1"`
	TotalRows  int                                                                `json:"totalRows" openapi:"example:1"`
	TotalPages int                                                                `json:"totalPages" openapi:"example:1"`
	Items      []CustomersQueryContractualMobilityRestrictionLogsResponseDataItem `json:"items" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsResponseDataItem"`
}

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsResponse
 */
type CustomersQueryContractualMobilityRestrictionLogsResponse struct {
	StatusCode int                                                          `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryContractualMobilityRestrictionLogsResponseData `json:"data" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsResponseData"`
}
