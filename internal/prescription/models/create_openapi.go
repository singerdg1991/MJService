package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/prescription/domain"
)

/*
 * @apiDefine: PrescriptionsCreateResponseData
 */
type PrescriptionsCreateResponseData struct {
	ID             uint                         `json:"id" openapi:"example:1"`
	CustomerID     uint                         `json:"customerId" openapi:"example:1;required"`
	Customer       *domain.PrescriptionCustomer `json:"customer" openapi:"$ref:PrescriptionCustomer;type:object;"`
	Title          string                       `json:"title" openapi:"example:John;required"`
	DateTime       string                       `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	DoctorFullName string                       `json:"doctorFullName" openapi:"example:John Doe;required"`
	StartDate      string                       `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate        string                       `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Status         string                       `json:"status" openapi:"example:active;required"`
	Attachments    []types.UploadMetadata       `json:"attachments" openapi:"$ref:UploadMetadata;type:array"`
	CreatedAt      string                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      string                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      string                       `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: PrescriptionsCreateResponse
 */
type PrescriptionsCreateResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200;"`
	Data       PrescriptionsCreateResponseData `json:"data" openapi:"$ref:PrescriptionsCreateResponseData"`
}
