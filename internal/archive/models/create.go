package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
)

/*
 * @apiDefine: ArchivesCreateRequestBody
 */
type ArchivesCreateRequestBody struct {
	Title       string              `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Subject     string              `json:"subject" openapi:"example:subject;required;maxLen:100;minLen:2;"`
	Description string              `json:"description" openapi:"example:description;required;maxLen:100;minLen:2;"`
	Date        string              `json:"date" openapi:"example:2021-01-01;required;"`
	Time        string              `json:"time" openapi:"example:00:00:00;required;"`
	DateTime    time.Time           `json:"-" openapi:"ignored"`
	UserID      uint                `json:"userId,string" openapi:"example:1;required;"`
	Attachments []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *ArchivesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":       govalidity.New("title").MinMaxLength(3, 25).Required(),
		"subject":     govalidity.New("subject").Required(),
		"description": govalidity.New("description").Optional(),
		"date":        govalidity.New("date").Required(),
		"time":        govalidity.New("time").Required(),
		"userId":      govalidity.New("userId").Int().Required(),
		"attachments": govalidity.New("attachments").Files(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate uploaded files size
	fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.ArchiveConfig.MaxUploadSizeLimit)
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
