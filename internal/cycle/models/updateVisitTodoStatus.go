package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/cycle/config"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: CyclesUpdateVisitTodoStatusRequestParams
 */
type CyclesUpdateVisitTodoStatusRequestParams struct {
	VisitTodoID int `json:"visittodoid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

// ValidateParams validates the parameters of the CyclesUpdateVisitTodoStatusRequestParams struct.
//
// It takes in a govalidity.Params object and checks if the "visittodoid" field is present and is of type int.
// If there are any validation errors, it returns a govalidity.ValidityResponseErrors object with the errors.
// Otherwise, it returns nil.
func (data *CyclesUpdateVisitTodoStatusRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"visittodoid": govalidity.New("visittodoid").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CyclesUpdateVisitTodoStatusRequestBody
 */
type CyclesUpdateVisitTodoStatusRequestBody struct {
	Status        string              `json:"status" openapi:"example:done"`
	Attachments   []*govalidityt.File `json:"attachments" openapi:"format:binary;type:array"`
	NotDoneReason *string             `json:"notDoneReason" openapi:"example:reason"`
}

func (data *CyclesUpdateVisitTodoStatusRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffTypes":    govalidity.New("staffTypes"),
		"attachments":   govalidity.New("attachments").Files(),
		"notDoneReason": govalidity.New("notDoneReason").Optional(),
	}

	// Check if the request body has error
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate uploaded files size
	fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.CycleConfig.MaxUploadSizeLimit)
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

	// Validate notDoneReason
	if data.Status == constants.TODO_STATUS_DONE {
		if data.NotDoneReason != nil {
			return govalidity.ValidityResponseErrors{
				"notDoneReason": []string{
					"Reason is not required for done status",
				},
			}
		}
	} else {
		if data.NotDoneReason == nil {
			return govalidity.ValidityResponseErrors{
				"notDoneReason": []string{
					"Reason is required for not done status",
				},
			}
		}
	}

	return nil
}
