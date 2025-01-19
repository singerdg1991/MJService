package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: PermissionsCreateRequestBody
 */
type PermissionsCreateRequestBody struct {
	Name  string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Title string `json:"title" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *PermissionsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":  govalidity.New("name").MinMaxLength(3, 25).Required(),
		"title": govalidity.New("title").MinMaxLength(3, 25).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
