package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersQueryMedicinesFilterType
 */
type CustomersQueryMedicinesFilterType struct {
}

/*
 * @apiDefine: CustomersQueryMedicinesSortValue
 */
type CustomersQueryMedicinesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryMedicinesSortType
 */
type CustomersQueryMedicinesSortType struct {
	CreatedAt CustomersQueryMedicinesSortValue `json:"created_at,omitempty" openapi:"$ref:CustomersQueryMedicinesSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryMedicinesRequestParams
 */
type CustomersQueryMedicinesRequestParams struct {
	ID         int                               `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                               `json:"customerId,string,omitempty" openapi:"example:1"`
	Page       int                               `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryMedicinesFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryMedicinesFilterType;in:query"`
	Sorts      CustomersQueryMedicinesSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryMedicinesSortType;in:query"`
}

func (data *CustomersQueryMedicinesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Optional(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters":    govalidity.Schema{},
		"sorts": govalidity.Schema{
			"created_at": govalidity.Schema{
				"op": govalidity.New("sort.created_at.op"),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
