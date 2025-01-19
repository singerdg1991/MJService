package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: CustomerServiceServiceType
 */
type CustomerServiceServiceType struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: CustomerServiceService
 */
type CustomerServiceService struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:John;required"`
}

/*
 * @apiDefine: CustomerServiceNurseWish
 */
type CustomerServiceNurseWish struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:John;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:John;required"`
}

/*
 * @apiDefine: CustomerServiceGrade
 */
type CustomerServiceGrade struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
	Grade       int    `json:"grade" openapi:"example:0;required"`
	Color       string `json:"color" openapi:"example:#000000;required"`
}

/*
 * @apiDefine: CustomerServiceCustomerUser
 */
type CustomerServiceCustomerUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:John;required"`
}

/*
 * @apiDefine: CustomerServiceCustomer
 */
type CustomerServiceCustomer struct {
	ID    uint                         `json:"id" openapi:"example:1"`
	KeyNo string                       `json:"keyNo" openapi:"example:John;required"`
	User  *CustomerServiceCustomerUser `json:"user" openapi:"$ref:CustomerServiceCustomerUser"`
}

/*
 * @apiDefine: CustomerServices
 */
type CustomerServices struct {
	ID                  uint                        `json:"id" openapi:"example:1"`
	CustomerID          uint                        `json:"customerId" openapi:"example:1"`
	Customer            *CustomerServiceCustomer    `json:"customer" openapi:"$ref:CustomerServiceCustomer"`
	ServiceID           *uint                       `json:"serviceId" openapi:"example:1"`
	Service             *CustomerServiceService     `json:"service" openapi:"$ref:CustomerServiceService"`
	ServiceTypeID       *uint                       `json:"serviceTypeId" openapi:"example:1"`
	ServiceType         *CustomerServiceServiceType `json:"serviceType" openapi:"$ref:CustomerServiceServiceType"`
	GradeID             *uint                       `json:"gradeId" openapi:"example:1"`
	Grade               *CustomerServiceGrade       `json:"grade" openapi:"$ref:CustomerServiceGrade"`
	NurseWishID         *uint                       `json:"staffWishId" openapi:"example:1"`
	NurseWish           *CustomerServiceNurseWish   `json:"staffWish" openapi:"$ref:CustomerServiceNurseWish"`
	ReportType          string                      `json:"reportType" openapi:"example:reportType"`
	TimeValue           *time.Time                  `json:"timeValue" openapi:"example:2021-01-01T00:00:00Z"`
	Repeat              *string                     `json:"repeat" openapi:"example:weekly"`
	VisitType           *string                     `json:"visitType" openapi:"example:online"`
	ServiceLengthMinute *uint                       `json:"serviceLengthMinute" openapi:"example:60"`
	StartTimeRange      *time.Time                  `json:"startTimeRange" openapi:"example:2021-01-01T00:00:00Z"`
	EndTimeRange        *time.Time                  `json:"endTimeRange" openapi:"example:2021-01-01T00:00:00Z"`
	Description         *string                     `json:"description" openapi:"example:John;required"`
	PaymentMethod       string                      `json:"paymentMethod" openapi:"example:own;required"`
	HomeCareFee         *uint                       `json:"homeCareFee" openapi:"example:100"`
	CityCouncilFee      *uint                       `json:"cityCouncilFee" openapi:"example:100"`
	CreatedAt           time.Time                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt           time.Time                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt           *time.Time                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerServices() *CustomerServices {
	return &CustomerServices{}
}

func (ns *CustomerServices) TableName() string {
	return "customerServices"
}

func (ns *CustomerServices) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
