package models

import (
	"encoding/json"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/customer/config"
	"net/http"
	"time"
)

/*
 * @apiDefine: CustomersUpdateMedicineRequestParams
 */
type CustomersUpdateMedicineRequestParams struct {
	ID         int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
	CustomerID int `json:"customerid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CustomersUpdateMedicineRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Required(),
		"customerid": govalidity.New("customerId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CustomersUpdateMedicinesRequestBodyHour
 */
type CustomersUpdateMedicinesRequestBodyHour struct {
	Hour        string `json:"hour" openapi:"example:12:00"`
	Description string `json:"description" openapi:"example:after meal"`
}

/*
 * @apiDefine: CustomersUpdateMedicinesRequestBody
 */
type CustomersUpdateMedicinesRequestBody struct {
	CustomerID                  int                                       `json:"customerId,string" openapi:"example:1"`
	PrescriptionID              uint                                      `json:"prescriptionId,string" openapi:"example:1"`
	MedicineID                  uint                                      `json:"medicineId,string" openapi:"example:1"`
	DosageAmount                uint                                      `json:"dosageAmount,string" openapi:"example:1"`
	DosageUnit                  string                                    `json:"dosageUnit" openapi:"example:gram"`
	Days                        *string                                   `json:"days" openapi:"example:[\"everyMonday\",\"everyTuesday\",\"everyWednesday\",\"everyThursday\",\"everyFriday\",\"everySaturday\",\"everySunday\"];type:array"`
	IsJustOneTime               string                                    `json:"isJustOneTime" openapi:"example:false"`
	Hours                       string                                    `json:"hours" openapi:"example:[{\"hour\":\"12:00\", \"description\":\"after meal\"}]"`
	StartDate                   string                                    `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate                     *string                                   `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Warning                     *string                                   `json:"warning" openapi:"example:warning"`
	IsUseAsNeeded               string                                    `json:"isUseAsNeeded" openapi:"example:false"`
	Attachments                 []*govalidityt.File                       `json:"attachments" openapi:"format:binary;type:array"`
	PreviousAttachments         string                                    `json:"previousAttachments" openapi:"example:[{\"fileName\": \"424e5ebcf1e4b4f11707315705332860929.png\", \"fileSize\": 44547, \"path\": \"/uploads/staff\"}]"`
	PreviousAttachmentsMetadata []types.UploadMetadata                    `json:"-" openapi:"ignored"`
	StartDateAsDate             *time.Time                                `json:"-" openapi:"ignored"`
	EndDateAsDate               *time.Time                                `json:"-" openapi:"ignored"`
	DaysAsArray                 []string                                  `json:"-" openapi:"ignored"`
	IsJustOneTimeAsBool         bool                                      `json:"-" openapi:"ignored"`
	IsUseAsNeededAsBool         bool                                      `json:"-" openapi:"ignored"`
	HoursMetadata               []CustomersCreateMedicinesRequestBodyHour `json:"-" openapi:"ignored"`
}

func (data *CustomersUpdateMedicinesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"customerId":     govalidity.New("customerId").Int().Required(),
		"prescriptionId": govalidity.New("prescriptionId").Int().Required(),
		"medicineId":     govalidity.New("medicineId").Int().Required(),
		"dosageAmount":   govalidity.New("dosageAmount").Int().Required(),
		"dosageUnit": govalidity.New("dosageUnit").In([]string{
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_GRAM,
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_MILLIGRAM,
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_MICROGRAM,
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_LITER,
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_MILLILITER,
			constants.CUSTOMER_MEDICINE_DOSAGE_UNIT_TEASPOON,
		}).Required(),
		"days":                govalidity.New("days"),
		"isJustOneTime":       govalidity.New("isJustOneTime").In([]string{"false", "true"}).Required(),
		"hours":               govalidity.New("hours"),
		"start_date":          govalidity.New("start_date").Required(),
		"end_date":            govalidity.New("end_date"),
		"warning":             govalidity.New("warning"),
		"isUseAsNeeded":       govalidity.New("isUseAsNeeded").In([]string{"false", "true"}).Required(),
		"attachments":         govalidity.New("attachments").Files(),
		"previousAttachments": govalidity.New("previousAttachments").Optional(),
	}

	// Validate bodyS
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Try to unmarshal Hours
	if data.Hours != "" {
		if err := json.Unmarshal([]byte(data.Hours), &data.HoursMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"hours": []string{"hours is not a valid JSON"},
			}
		}
	}

	// Try to unmarshal PreviousAttachmentsMetadata
	if data.PreviousAttachments != "" {
		if err := json.Unmarshal([]byte(data.PreviousAttachments), &data.PreviousAttachmentsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{"previousAttachments is not a valid JSON"},
			}
		}
	}

	// Validate start_date
	startDate, err := time.Parse(time.RFC3339, data.StartDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"Invalid date format"},
		}
	}
	data.StartDateAsDate = &startDate

	// Validate end_date
	if data.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *data.EndDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"Invalid date format"},
			}
		}
		data.EndDateAsDate = &endDate
	}

	// Validate days
	if data.Days != nil {
		// Try to unmarshal Days
		var days interface{}
		if err := json.Unmarshal([]byte(*data.Days), &days); err != nil {
			return govalidity.ValidityResponseErrors{
				"days": []string{"days is not a valid JSON"},
			}
		}
		daysArray, ok := days.([]interface{})
		if !ok {
			return govalidity.ValidityResponseErrors{
				"days": []string{"days is not a valid JSON"},
			}
		}
		for _, day := range daysArray {
			day, ok := day.(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"days": []string{"days is not a valid JSON"},
				}
			}
			data.DaysAsArray = append(data.DaysAsArray, day)
		}
	}

	// Validate hours
	for _, hour := range data.HoursMetadata {
		// Validate hour
		_, err := time.Parse("15:04", hour.Hour)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"hours": []string{"Invalid hours format"},
			}
		}
	}

	// Validate isJustOneTime
	if data.IsJustOneTime == "true" {
		data.IsJustOneTimeAsBool = true
	} else {
		data.IsJustOneTimeAsBool = false
	}

	// Validate isUseAsNeeded
	if data.IsUseAsNeeded == "true" {
		data.IsUseAsNeededAsBool = true
	} else {
		data.IsUseAsNeededAsBool = false
	}

	// Validate uploaded files size
	fileErrs := sharedutils.ValidateUploadFilesSize("attachments", data.Attachments, config.CustomerConfig.MaxUploadSizeLimit)
	if fileErrs != nil {
		return fileErrs
	}

	// Validate uploaded files mime type
	fileErrs = sharedutils.ValidateUploadFilesMimeType("attachments", data.Attachments, []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
	})
	if fileErrs != nil {
		return fileErrs
	}

	return nil
}
