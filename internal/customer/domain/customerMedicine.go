package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: CustomerMedicineHour
 */
type CustomerMedicineHour struct {
	Hour        string `json:"hour" openapi:"example:08:00"`
	Description string `json:"description" openapi:"example:description;required"`
}

/*
 * @apiDefine: CustomerMedicinePrescription
 */
type CustomerMedicinePrescription struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
}

/*
 * @apiDefine: CustomerMedicineMedicine
 */
type CustomerMedicineMedicine struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:name;required"`
}

/*
 * @apiDefine: CustomerMedicine
 */
type CustomerMedicine struct {
	ID             uint                          `json:"id" openapi:"example:1"`
	CustomerID     uint                          `json:"customerId" openapi:"example:1"`
	PrescriptionID uint                          `json:"prescriptionId" openapi:"example:1"`
	Prescription   *CustomerMedicinePrescription `json:"prescription" openapi:"$ref:CustomerMedicinePrescription"`
	MedicineID     uint                          `json:"medicineId" openapi:"example:1"`
	Medicine       *CustomerMedicineMedicine     `json:"medicine" openapi:"$ref:CustomerMedicineMedicine"`
	DosageAmount   uint                          `json:"dosageAmount" openapi:"example:1"`
	DosageUnit     string                        `json:"dosageUnit" openapi:"example:gram"`
	Days           []string                      `json:"days" openapi:"example:[\"everyMonday\",\"everyTuesday\",\"everyWednesday\",\"everyThursday\",\"everyFriday\",\"everySaturday\",\"everySunday\"];type:array"`
	IsJustOneTime  bool                          `json:"isJustOneTime" openapi:"example:false"`
	Hours          []CustomerMedicineHour        `json:"hours" openapi:"$ref:CustomerMedicineHour;type:array"`
	StartDate      *time.Time                    `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate        *time.Time                    `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Warning        *string                       `json:"warning" openapi:"example:warning"`
	IsUseAsNeeded  bool                          `json:"isUseAsNeeded" openapi:"example:false"`
	Attachments    []*types.UploadMetadata       `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	CreatedAt      time.Time                     `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time                     `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time                    `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomerMedicine() *CustomerMedicine {
	return &CustomerMedicine{}
}

func (ns *CustomerMedicine) TableName() string {
	return "customerMedicines"
}

func (ns *CustomerMedicine) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
