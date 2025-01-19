package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ServiceOptionsDeleteRequestBody
 */
type ServiceOptionsDeleteRequestBody struct {
	IDs      interface{} `json:"ids" openapi:"example:[1,2,3];type:array;required;"`
	IDsInt64 []int64     `json:"-" openapi:"ignored"`
}

func (data *ServiceOptionsDeleteRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"ids": govalidity.New("ids"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
