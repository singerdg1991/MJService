package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"net/http"
)

/*
 * @apiDefine: StaffsUpdateLicensesRequestParams
 */
type StaffsUpdateLicensesRequestParams struct {
	StaffID int `json:"staffId,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *StaffsUpdateLicensesRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId": govalidity.New("staffId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: StaffsCreateLicensesRequestBodyLicense
 */
type StaffsCreateLicensesRequestBodyLicense struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateLicensesRequestBody
 */
type StaffsCreateLicensesRequestBody struct {
	StaffID     int                                    `json:"staffId,string" openapi:"example:1"`
	License     StaffsCreateLicensesRequestBodyLicense `json:"license" openapi:"$ref:StaffsCreateLicensesRequestBodyLicense"`
	ExpireDate  string                                 `json:"expire_date" openapi:"example:2020-01-01T00:00:00Z"`
	Attachments []*govalidityt.File                    `json:"attachments" openapi:"format:binary;type:array"`
}

func (data *StaffsCreateLicensesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId":     govalidity.New("staffId").Int().Required(),
		"license":     govalidity.New("license").Optional(),
		"expire_date": govalidity.New("expire_date").Required(),
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
