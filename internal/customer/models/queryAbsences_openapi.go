package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersQueryAbsencesResponseDataItem
 */
type CustomersQueryAbsencesResponseDataItem struct {
	ID          int                            `json:"id" openapi:"example:1"`
	CustomerID  int                            `json:"customerId" openapi:"example:1"`
	Customer    domain.CustomerAbsenceCustomer `json:"customer" openapi:"$ref:CustomerAbsenceCustomer"`
	StartDate   string                         `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate     *string                        `json:"end_date" openapi:"example:2020-01-01T00:00:00Z;nullable"`
	Reason      *string                        `json:"reason" openapi:"example:reason;nullable"`
	Attachments []types.UploadMetadata         `json:"attachments" openapi:"$ref:UploadMetadata;type:array"`
	CreatedAt   string                         `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt   string                         `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt   *string                        `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z;nullable"`
}

/*
 * @apiDefine: CustomersQueryAbsencesResponseData
 */
type CustomersQueryAbsencesResponseData struct {
	Limit      int                                      `json:"limit" openapi:"example:10"`
	Offset     int                                      `json:"offset" openapi:"example:0"`
	Page       int                                      `json:"page" openapi:"example:1"`
	TotalRows  int                                      `json:"totalRows" openapi:"example:1"`
	TotalPages int                                      `json:"totalPages" openapi:"example:1"`
	Items      []CustomersQueryAbsencesResponseDataItem `json:"items" openapi:"$ref:CustomersQueryAbsencesResponseDataItem"`
}

/*
 * @apiDefine: CustomersQueryAbsencesResponse
 */
type CustomersQueryAbsencesResponse struct {
	StatusCode int                                `json:"statusCode" openapi:"example:200"`
	Data       CustomersQueryAbsencesResponseData `json:"data" openapi:"$ref:CustomersQueryAbsencesResponseData"`
}
