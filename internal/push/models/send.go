package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: PushesSendRequestBody
 */
type PushesSendRequestBody struct {
	UserID int    `json:"userId" openapi:"example:1;required;"`
	Title  string `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Body   string `json:"body" openapi:"example:body;required;maxLen:100;minLen:2;"`
}

func (data *PushesSendRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userId": govalidity.New("userId").Int().Min(1).Required(),
		"title":  govalidity.New("title").MinMaxLength(2, 100).Required(),
		"body":   govalidity.New("body").MinMaxLength(2, 100).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
