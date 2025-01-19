package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: AIsCreateRequestBody
 */
type AIsCreateRequestBody struct {
	Subject string `json:"subject" openapi:"example:subject;required;maxLen:100;minLen:2;"`
}

// ValidateBody validates the body of an HTTP request against a predefined schema.
//
// The function takes an HTTP request as a parameter and checks its body against a schema defined for the AIsCreateRequestBody struct.
// It returns a govalidity.ValidityResponseErrors object containing any validation errors that occurred.
func (data *AIsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"subject": govalidity.New("subject").MaxLength(1000).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
