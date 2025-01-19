package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/ticket/config"
	"net/http"
)

/*
 * @apiDefine: TicketsCreateMessageRequestParams
 */
type TicketsCreateMessageRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *TicketsCreateMessageRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: TicketsCreateMessageRequestBody
 */
type TicketsCreateMessageRequestBody struct {
	Message     string              `json:"message" openapi:"example:message;required;maxLen:100;minLen:2;"`
	Attachments []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *TicketsCreateMessageRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"recipientId": govalidity.New("recipientId"),
		"message":     govalidity.New("message").Required(),
		"attachments": govalidity.New("attachments").Files(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate attachments
	if len(data.Attachments) > 0 {
		// Validate uploaded files size
		fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.TicketConfig.MaxUploadSizeLimit)
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
