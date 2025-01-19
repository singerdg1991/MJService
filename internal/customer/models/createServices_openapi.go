package models

import (
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreateServicesResponseData
 */
type CustomersCreateServicesResponseData struct {
	ID                  int                                `json:"id" openapi:"example:1"`
	CustomerID          uint                               `json:"customerId" openapi:"example:1"`
	Customer            *domain.CustomerServiceCustomer    `json:"customer" openapi:"$ref:CustomerServiceCustomer"`
	ServiceID           *uint                              `json:"serviceId" openapi:"example:1"`
	Service             *domain.CustomerServiceService     `json:"service" openapi:"$ref:CustomerServiceService"`
	ServiceTypeID       *uint                              `json:"serviceTypeId" openapi:"example:1"`
	ServiceType         *domain.CustomerServiceServiceType `json:"serviceType" openapi:"$ref:CustomerServiceServiceType"`
	GradeID             *uint                              `json:"gradeId" openapi:"example:1"`
	Grade               *domain.CustomerServiceGrade       `json:"grade" openapi:"$ref:CustomerServiceGrade"`
	NurseWishID         *uint                              `json:"staffWishId" openapi:"example:1"`
	NurseWish           *domain.CustomerServiceNurseWish   `json:"staffWish" openapi:"$ref:CustomerServiceNurseWish"`
	ReportType          string                             `json:"reportType" openapi:"example:reportType"`
	TimeValue           string                             `json:"timeValue" openapi:"example:2021-01-01T00:00:00Z"`
	Repeat              *string                            `json:"repeat" openapi:"example:weekly"`
	VisitType           *string                            `json:"visitType" openapi:"example:online"`
	ServiceLengthMinute *uint                              `json:"serviceLengthMinute" openapi:"example:60"`
	StartTimeRange      string                             `json:"startTimeRange" openapi:"example:2021-01-01T00:00:00Z"`
	EndTimeRange        string                             `json:"endTimeRange" openapi:"example:2021-01-01T00:00:00Z"`
	Description         *string                            `json:"description" openapi:"example:John;required"`
	CreatedAt           string                             `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt           string                             `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt           *string                            `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z;nullable"`
}

/*
 * @apiDefine: CustomersCreateServicesResponse
 */
type CustomersCreateServicesResponse struct {
	StatusCode int                                 `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateServicesResponseData `json:"data" openapi:"$ref:CustomersCreateServicesResponseData"`
}
