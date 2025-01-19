package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CustomersQueryServicesFilterType
 */
type CustomersQueryServicesFilterType struct {
}

/*
 * @apiDefine: CustomersQueryServicesSortValue
 */
type CustomersQueryServicesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryServicesSortType
 */
type CustomersQueryServicesSortType struct {
}

/*
 * @apiDefine: CustomersQueryServicesRequestParams
 */
type CustomersQueryServicesRequestParams struct {
	ID         int                              `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                              `json:"customerId,string,omitempty" openapi:"example:1;required"`
	Page       int                              `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                              `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryServicesFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryServicesFilterType;in:query"`
	Sorts      CustomersQueryServicesSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryServicesSortType;in:query"`
}

func (data *CustomersQueryServicesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Required(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters":    govalidity.Schema{},
		"sorts":      govalidity.Schema{},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
