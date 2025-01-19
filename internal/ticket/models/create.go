package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/ticket/config"
	"net/http"
)

/*
 * @apiDefine: TicketsCreateRequestBody
 */
type TicketsCreateRequestBody struct {
	Title        string              `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Description  string              `json:"description" openapi:"example:description;required;maxLen:100;minLen:2;"`
	Priority     string              `json:"priority" openapi:"example:low;required;maxLen:100;minLen:2;"`
	UserID       *uint               `json:"userId" openapi:"example:1;required;"`
	DepartmentID *uint               `json:"departmentId" openapi:"example:1;required;"`
	Attachments  []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *TicketsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":        govalidity.New("title").Required(),
		"description":  govalidity.New("description").Required(),
		"priority":     govalidity.New("priority").Required(),
		"userId":       govalidity.New("userId"),
		"departmentId": govalidity.New("departmentId"),
		"attachments":  govalidity.New("attachments").Files(),
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

	// Check only userId or departmentId is provided
	if data.UserID != nil && data.DepartmentID != nil {
		return govalidity.ValidityResponseErrors{
			"userId": []string{"Only one of userId or departmentId can be provided"},
		}
	}

	// Check userID is integer value and greater than 0
	if data.UserID != nil && *data.UserID <= 0 {
		return govalidity.ValidityResponseErrors{
			"userId": []string{"userId must be integer value and greater than 0"},
		}
	}

	// Check departmentID is integer value and greater than 0
	if data.DepartmentID != nil && *data.DepartmentID <= 0 {
		return govalidity.ValidityResponseErrors{
			"departmentId": []string{"departmentId must be integer value and greater than 0"},
		}
	}

	return nil
}
