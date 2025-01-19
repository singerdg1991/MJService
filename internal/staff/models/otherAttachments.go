package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"net/http"
)

/*
 * @apiDefine: StaffsCreateOtherAttachmentsRequestBody
 */
type StaffsCreateOtherAttachmentsRequestBody struct {
	StaffID     int                 `json:"staffId,string" openapi:"example:1"`
	Title       string              `json:"title" openapi:"example:attachment title"`
	Attachments []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *StaffsCreateOtherAttachmentsRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId":     govalidity.New("staffId").Int().Required(),
		"title":       govalidity.New("title").Required(),
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
