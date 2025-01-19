package models

import "github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"

/*
 * @apiDefine: HolidaysResponseData
 */
type HolidaysResponseData struct {
	ID             uint               `json:"id" openapi:"example:1"`
	StartDate      string             `json:"start_date" openapi:"example:2020-01-01"`
	EndDate        string             `json:"end_date" openapi:"example:2020-01-01"`
	Title          string             `json:"title" openapi:"example:title"`
	PaymentType    string             `json:"paymentType" openapi:"example:withSalary"`
	Description    *string            `json:"description" openapi:"example:description"`
	Status         string             `json:"status" openapi:"example:pending"`
	RejectedReason string             `json:"rejectedReason" openapi:"example:rejectedReason"`
	AcceptedAt     string             `json:"accepted_at" openapi:"example:2021-01-01T00:00:00Z"`
	RejectedAt     string             `json:"rejected_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt      string             `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      string             `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      string             `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedBy      domain.HolidayUser `json:"createdBy" openapi:"$ref:HolidayUser"`
	UpdatedBy      domain.HolidayUser `json:"updatedBy" openapi:"$ref:HolidayUser"`
}

/*
 * @apiDefine: HolidaysCreateResponse
 */
type HolidaysCreateResponse struct {
	StatusCode int                  `json:"statusCode" openapi:"example:200"`
	Data       HolidaysResponseData `json:"data" openapi:"$ref:HolidaysResponseData"`
}
