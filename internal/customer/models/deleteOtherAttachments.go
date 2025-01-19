package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersDeleteOtherAttachmentsRequestBody
 */
type CustomersDeleteOtherAttachmentsRequestBody struct {
	IDs      interface{} `json:"ids" openapi:"example:[1,2,3];type:array;required;"`
	IDsInt64 []int64     `json:"-" openapi:"ignored"`
	UserID   int         `json:"userId" openapi:"example:1;required;"`
}

func (data *CustomersDeleteOtherAttachmentsRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"ids":    govalidity.New("ids"),
		"userId": govalidity.New("userId").Int().Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
