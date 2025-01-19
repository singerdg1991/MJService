package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ServicesCreateRequestBody
 */
type ServicesCreateRequestBody struct {
	Name        string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description string `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *ServicesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":        govalidity.New("name").MinMaxLength(3, 25).Required(),
		"description": govalidity.New("description").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
