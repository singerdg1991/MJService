package models

import (
	"github.com/hoitek/Govalidity/govalidityt"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/prescription/config"
	"github.com/hoitek/Maja-Service/internal/prescription/constants"
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: PrescriptionsCreateRequestBody
 */
type PrescriptionsCreateRequestBody struct {
	CustomerID      int                 `json:"customerId,string" openapi:"example:1"`
	Title           string              `json:"title" openapi:"example:Title"`
	DateTime        string              `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	DoctorFullName  string              `json:"doctorFullName" openapi:"example:John Doe"`
	StartDate       string              `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate         string              `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Status          string              `json:"status" openapi:"example:active"`
	Attachments     []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
	DateTimeAsDate  *time.Time          `json:"-" openapi:"ignored"`
	StartDateAsDate *time.Time          `json:"-" openapi:"ignored"`
	EndDateAsDate   *time.Time          `json:"-" openapi:"ignored"`
}

func (data *PrescriptionsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"customerId":     govalidity.New("customerId").Int().Required(),
		"title":          govalidity.New("title").MinLength(3).Required(),
		"datetime":       govalidity.New("datetime").Required(),
		"doctorFullName": govalidity.New("doctorFullName").MinLength(3).Required(),
		"start_date":     govalidity.New("start_date").Required(),
		"end_date":       govalidity.New("end_date"),
		"status": govalidity.New("status").In([]string{
			constants.PRESCRIPTION_STATUS_ACTIVE,
			constants.PRESCRIPTION_STATUS_FILED,
			constants.PRESCRIPTION_STATUS_EXPIRED,
		}),
		"attachments": govalidity.New("attachments").Files(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	if data.DateTime != "" {
		datetime, err := time.Parse(time.RFC3339, data.DateTime)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"datetime": []string{"Invalid datetime format"},
			}
		}
		data.DateTimeAsDate = &datetime
	}

	// Validate start_date
	if data.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, data.StartDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"start_date": []string{"Invalid start_date format"},
			}
		}
		data.StartDateAsDate = &startDate
	}

	// Validate end_date
	if data.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, data.EndDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"Invalid end_date format"},
			}
		}
		data.EndDateAsDate = &endDate
	}

	// Validate uploaded files size
	fileErrs := sharedutils.ValidateUploadFilesSize("attachments", data.Attachments, config.PrescriptionConfig.MaxUploadSizeLimit)
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
