package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: DiagnosesCreateRequestBody
 */
type DiagnosesCreateRequestBody struct {
	Title       string  `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Code        string  `json:"code" openapi:"example:code;required;maxLen:100;minLen:2;"`
	Description *string `json:"description" openapi:"example:description;maxLen:100;minLen:2;"`
}

func (data *DiagnosesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":       govalidity.New("title").MinMaxLength(3, 25).Required(),
		"code":        govalidity.New("code").MinMaxLength(3, 25).Required(),
		"description": govalidity.New("description").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
