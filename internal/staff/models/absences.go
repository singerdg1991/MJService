package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"net/http"
)

/*
 * @apiDefine: StaffsCreateAbsencesRequestBody
 */
type StaffsCreateAbsencesRequestBody struct {
	StaffID     int                 `json:"staffId,string" openapi:"example:1"`
	StartDate   string              `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate     string              `json:"end_date" openapi:"example:2020-01-01T00:00:00Z"`
	Reason      string              `json:"reason" openapi:"example:reason"`
	Attachments []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *StaffsCreateAbsencesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId":     govalidity.New("staffId").Int().Required(),
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

	// Validate uploaded files size
	fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.StaffConfig.MaxUploadSizeLimit)
	if fileErrs != nil {
		return fileErrs
	}

	// Validate uploaded files mime type
	fileErrs = utils.ValidateUploadFilesMimeType("attachments", data.Attachments, []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
	})
	if fileErrs != nil {
		return fileErrs
	}

	return nil
}
