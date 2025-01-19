package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsDeleteRequestBody
 */
type StaffsDeleteRequestBody struct {
	IDs string `json:"ids" openapi:"example:1,2,3;required;"`
}

func (data *StaffsDeleteRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"ids": govalidity.New("ids").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
