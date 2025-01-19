package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

type StaffLicenses struct {
	ID          int64                   `json:"id"`
	StaffID     int64                   `json:"staffId"`
	LicenseID   int64                   `json:"licenseId"`
	ExpireDate  *time.Time              `json:"expire_date" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt   time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffLicenses() *StaffLicenses {
	return &StaffLicenses{}
}

func (ns *StaffLicenses) TableName() string {
	return "staffLicenses"
}

func (ns *StaffLicenses) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
