package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: PrescriptionCustomer
 */
type PrescriptionCustomer struct {
	ID        int64  `json:"id" openapi:"example:1"`
	UserID    int64  `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: Prescription
 */
type Prescription struct {
	ID             uint                    `json:"id" openapi:"example:1"`
	CustomerID     uint                    `json:"customerId" openapi:"example:1;required"`
	Customer       *PrescriptionCustomer   `json:"customer" openapi:"$ref:PrescriptionCustomer;type:object;"`
	Title          string                  `json:"title" openapi:"example:John;required"`
	DateTime       *time.Time              `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	DoctorFullName string                  `json:"doctorFullName" openapi:"example:John Doe;required"`
	StartDate      *time.Time              `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate        *time.Time              `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Status         string                  `json:"status" openapi:"example:active;required"`
	Attachments    []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt      time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Prescription) TableName() string {
	return "prescriptions"
}

func (u *Prescription) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
