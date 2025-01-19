package models

import (
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CustomersUpdatePersonalInfoRequestParams
 */
type CustomersUpdatePersonalInfoRequestParams struct {
	ID         int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
	CustomerID int `json:"customerid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CustomersUpdatePersonalInfoRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Required(),
		"customerid": govalidity.New("customerId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
