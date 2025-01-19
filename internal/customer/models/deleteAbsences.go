package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersDeleteAbsencesRequestBody
 */
type CustomersDeleteAbsencesRequestBody struct {
	IDs        interface{} `json:"ids" openapi:"example:[1,2,3];type:array;required;"`
	IDsInt64   []int64     `json:"-" openapi:"ignored"`
	CustomerID int         `json:"customerId" openapi:"example:1;required;"`
}

func (data *CustomersDeleteAbsencesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"ids":        govalidity.New("ids"),
		"customerId": govalidity.New("customerId").Int().Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
