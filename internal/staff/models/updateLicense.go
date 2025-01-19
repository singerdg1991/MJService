package models

import (
	"encoding/json"
	"fmt"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"net/http"
)

/*
 * @apiDefine: StaffsUpdateLicenseRequestParams
 */
type StaffsUpdateLicenseRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *StaffsUpdateLicenseRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: StaffsUpdateLicenseRequestBody
 */
type StaffsUpdateLicenseRequestBody struct {
	License                     StaffsCreateLicensesRequestBodyLicense `json:"license" openapi:"$ref:StaffsCreateLicensesRequestBodyLicense"`
	ExpireDate                  string                                 `json:"expire_date" openapi:"example:2020-01-01T00:00:00Z"`
	Attachments                 []*govalidityt.File                    `json:"attachments" openapi:"format:binary;type:array"`
	PreviousAttachments         string                                 `json:"previousAttachments" openapi:"example:[{\"fileName\": \"424e5ebcf1e4b4f11707315705332860929.png\", \"fileSize\": 44547, \"path\": \"/uploads/staff\"}]"`
	PreviousAttachmentsMetadata []types.UploadMetadata                 `json:"-" openapi:"ignored"`
}

func (data *StaffsUpdateLicenseRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"license":             govalidity.New("license").Optional(),
		"expire_date":         govalidity.New("expire_date").Required(),
		"attachments":         govalidity.New("attachments").Files(),
		"previousAttachments": govalidity.New("previousAttachments").Optional(),
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

	// Try to unmarshal PreviousAttachmentsMetadata
	if data.PreviousAttachments != "" {
		if err := json.Unmarshal([]byte(data.PreviousAttachments), &data.PreviousAttachmentsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{"previousAttachments is not a valid JSON"},
			}
		}
	}

	// Validate PreviousAttachments
	for i, attachment := range data.PreviousAttachmentsMetadata {
		if attachment.FileName == "" {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{fmt.Sprintf("previousAttachments[%d].fileName is required", i)},
			}
		}
		if attachment.Path == "" {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{fmt.Sprintf("previousAttachments[%d].path is required", i)},
			}
		}
	}

	return nil
}
