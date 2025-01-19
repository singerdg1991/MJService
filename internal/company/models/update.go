package models

import (
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CompaniesUpdateRequestParams
 */
type CompaniesUpdateRequestParams struct {
	Name string `json:"name" openapi:"example:name;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CompaniesUpdateRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name": govalidity.New("name").MinMaxLength(3, 25).Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
