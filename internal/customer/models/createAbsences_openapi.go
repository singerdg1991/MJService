package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreateAbsencesResponseData
 */
type CustomersCreateAbsencesResponseData struct {
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
 * @apiDefine: CustomersCreateAbsencesResponse
 */
type CustomersCreateAbsencesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateAbsencesResponseData `json:"data" openapi:"$ref:CustomersCreateAbsencesResponseData"`
}
