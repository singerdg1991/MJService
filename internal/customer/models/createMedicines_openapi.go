package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreateMedicinesResponseData
 */
type CustomersCreateMedicinesResponseData struct {
	ID             uint                                 `json:"id" openapi:"example:1"`
	CustomerID     uint                                 `json:"customerId" openapi:"example:1"`
	PrescriptionID uint                                 `json:"prescriptionId" openapi:"example:1"`
	Prescription   *domain.CustomerMedicinePrescription `json:"prescription" openapi:"$ref:CustomerMedicinePrescription"`
	MedicineID     uint                                 `json:"medicineId" openapi:"example:1"`
	Medicine       *domain.CustomerMedicineMedicine     `json:"medicine" openapi:"$ref:CustomerMedicineMedicine"`
	DosageAmount   uint                                 `json:"dosageAmount" openapi:"example:1"`
	DosageUnit     string                               `json:"dosageUnit" openapi:"example:gram"`
	Days           interface{}                          `json:"days" openapi:"example:[\"everyMonday\",\"everyTuesday\",\"everyWednesday\",\"everyThursday\",\"everyFriday\",\"everySaturday\",\"everySunday\"];type:array"`
	IsJustOneTime  bool                                 `json:"isJustOneTime" openapi:"example:false"`
	Hours          []domain.CustomerMedicineHour        `json:"hours" openapi:"$ref:CustomerMedicineHour;type:array"`
	StartDate      *string                              `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate        *string                              `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Warning        *string                              `json:"warning" openapi:"example:warning"`
	IsUseAsNeeded  bool                                 `json:"isUseAsNeeded" openapi:"example:false"`
	Attachments    []types.UploadMetadata               `json:"attachments" openapi:"$ref:UploadMetadata;type:array"`
	CreatedAt      string                               `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt      string                               `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt      *string                              `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z;nullable"`
}

/*
 * @apiDefine: CustomersCreateMedicinesResponse
 */
type CustomersCreateMedicinesResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreateMedicinesResponseData `json:"data" openapi:"$ref:CustomersCreateMedicinesResponseData"`
}
