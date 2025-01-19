package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"github.com/hoitek/Maja-Service/utils"
	"net/http"
	"time"
)

/*
 * @apiDefine: CustomersCreateAbsencesRequestBody
 */
type CustomersCreateAbsencesRequestBody struct {
	CustomerID      int                 `json:"customerId,string" openapi:"example:1"`
	StartDate       string              `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate         *string             `json:"end_date" openapi:"example:2020-01-01T00:00:00Z"`
	Reason          string              `json:"reason" openapi:"example:reason"`
	Attachments     []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
	StartDateAsDate *time.Time          `json:"-" openapi:"ignored"`
	EndDateAsDate   *time.Time          `json:"-" openapi:"ignored"`
}

func (data *CustomersCreateAbsencesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"customerId":  govalidity.New("customerId").Int().Required(),
		"start_date":  govalidity.New("start_date").Required(),
		"end_date":    govalidity.New("end_date").Optional(),
		"reason":      govalidity.New("reason").Optional(),
		"attachments": govalidity.New("attachments").Files(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate startDate
	startDate, err := utils.TryParseToDateTime(data.StartDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"startDate is invalid"},
		}
	}
	data.StartDateAsDate = &startDate

	// Validate endDate
	if data.EndDate != nil {
		endDate, err := utils.TryParseToDateTime(*data.EndDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"endDate is invalid"},
			}
		}
		data.EndDateAsDate = &endDate
	}

	// Validate uploaded files size
	fileErrs := sharedutils.ValidateUploadFilesSize("attachments", data.Attachments, config.StaffConfig.MaxUploadSizeLimit)
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
