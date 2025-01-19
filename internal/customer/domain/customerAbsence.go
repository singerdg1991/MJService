package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: CustomerAbsenceCustomer
 */
type CustomerAbsenceCustomer struct {
	ID        int64  `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
	Email     string `json:"email" openapi:"example:sgh370@yahoo.com"`
}

/*
 * @apiDefine: CustomerAbsence
 */
type CustomerAbsence struct {
	ID          int64                    `json:"id" openapi:"example:1"`
	CustomerID  int64                    `json:"customerId" openapi:"example:1"`
	Customer    *CustomerAbsenceCustomer `json:"customer" openapi:"$ref:CustomerAbsenceCustomer"`
	StartDate   time.Time                `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate     *time.Time               `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Reason      *string                  `json:"reason" openapi:"example:reason"`
	Attachments []*types.UploadMetadata  `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt   time.Time                `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time                `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time               `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerAbsence() *CustomerAbsence {
	return &CustomerAbsence{}
}

func (ns *CustomerAbsence) TableName() string {
	return "customerAbsences"
}

func (ns *CustomerAbsence) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
