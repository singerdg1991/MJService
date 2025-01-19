package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsDeleteLicensesRequestBody
 */
type StaffsDeleteLicensesRequestBody struct {
	IDs     string `json:"ids" openapi:"example:1,2,3;required;"`
	StaffID int    `json:"staffId" openapi:"example:1;required;"`
}

func (data *StaffsDeleteLicensesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"ids":     govalidity.New("ids").Required(),
		"staffId": govalidity.New("staffId").Int().Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
