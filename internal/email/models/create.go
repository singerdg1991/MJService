package models

import (
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/email/config"
	"github.com/hoitek/Maja-Service/internal/email/constants"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: EmailsCreateRequestBody
 */
type EmailsCreateRequestBody struct {
	Email             string                          `json:"email" openapi:"example:sgh370@yahoo.com;required"`
	Title             string                          `json:"title" openapi:"example:title;required"`
	Subject           string                          `json:"subject" openapi:"example:subject;required"`
	Message           string                          `json:"message" openapi:"example:message;required"`
	Attachments       []*govalidityt.File             `json:"attachments" openapi:"format:binary;type:array"`
	Category          string                          `json:"category" openapi:"example:outbox;required"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *EmailsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"email":       govalidity.New("email").Email().Required(),
		"title":       govalidity.New("title").Required(),
		"subject":     govalidity.New("subject").Required(),
		"message":     govalidity.New("message").Required(),
		"attachments": govalidity.New("attachments").Files(),
		"category": govalidity.New("category").Required().In([]string{
			constants.EMAIL_CATEGORY_OUTBOX,
			constants.EMAIL_CATEGORY_DRAFT,
			constants.EMAIL_CATEGORY_ARCHIVE,
			constants.EMAIL_CATEGORY_TRASH,
		}),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate attachments
	if len(data.Attachments) > 0 {
		// Validate uploaded files size
		fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.EmailConfig.MaxUploadSizeLimit)
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
	}

	return nil
}
