package models

import (
	"encoding/json"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	sharedutils "github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"github.com/hoitek/Maja-Service/utils"
	"net/http"
	"time"
)

/*
 * @apiDefine: CustomersUpdateAbsenceRequestParams
 */
type CustomersUpdateAbsenceRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CustomersUpdateAbsenceRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id": govalidity.New("id").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CustomersUpdateAbsenceRequestBody
 */
type CustomersUpdateAbsenceRequestBody struct {
	StartDate                   string                 `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate                     *string                `json:"end_date" openapi:"example:2020-01-01T00:00:00Z"`
	Reason                      string                 `json:"reason" openapi:"example:reason"`
	Attachments                 []*govalidityt.File    `json:"attachments" openapi:"format:binary;type:array"`
	PreviousAttachments         string                 `json:"previousAttachments" openapi:"example:[{\"fileName\": \"424e5ebcf1e4b4f11707315705332860929.png\", \"fileSize\": 44547, \"path\": \"/uploads/staff\"}]"`
	PreviousAttachmentsMetadata []types.UploadMetadata `json:"-" openapi:"ignored"`
	StartDateAsDate             *time.Time             `json:"-" openapi:"ignored"`
	EndDateAsDate               *time.Time             `json:"-" openapi:"ignored"`
}

func (data *CustomersUpdateAbsenceRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"start_date":          govalidity.New("start_date").Required(),
		"end_date":            govalidity.New("end_date").Optional(),
		"reason":              govalidity.New("reason").Optional(),
		"attachments":         govalidity.New("attachments").Files(),
		"previousAttachments": govalidity.New("previousAttachments").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Try to unmarshal PreviousAttachmentsMetadata
	if data.PreviousAttachments != "" {
		if err := json.Unmarshal([]byte(data.PreviousAttachments), &data.PreviousAttachmentsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{"previousAttachments is not a valid JSON"},
			}
		}
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
