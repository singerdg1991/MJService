package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: RolesCreateRequestBody
 */
type RolesCreateRequestBody struct {
	Name             string      `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Permissions      interface{} `json:"permissions" openapi:"example:[1,2,3];type:array;required;"`
	PermissionsInt64 []int64     `json:"-" openapi:"ignored"`
}

func (data *RolesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":        govalidity.New("name").MinMaxLength(3, 25).Required(),
		"permissions": govalidity.New("permissions"),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
