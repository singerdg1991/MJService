package models

import (
	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomerCreditDetailsFilterType
 */
type CustomerCreditDetailsFilterType struct {
	BankAccountNumber  filters.FilterValue[string] `json:"bankAccountNumber,omitempty" openapi:"$ref:FilterValueString;example:{\"bankAccountNumber\":{\"op\":\"equals\",\"value\":\"1234567890\"}"`
	BillingAddressName filters.FilterValue[string] `json:"billingAddressName,omitempty" openapi:"$ref:FilterValueString;example:{\"billingAddressName\":{\"op\":\"equals\",\"value\":\"Home\"}"`
}

/*
 * @apiDefine: CustomerCreditDetailsSortValue
 */
type CustomerCreditDetailsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomerCreditDetailsSortType
 */
type CustomerCreditDetailsSortType struct {
	CreatedAt CustomerCreditDetailsSortValue `json:"created_at,omitempty" openapi:"$ref:CustomerCreditDetailsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryCreditDetailsRequestParams
 */
type CustomersQueryCreditDetailsRequestParams struct {
	ID         int                             `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                             `json:"customerId,string,omitempty" openapi:"example:1"`
	Page       int                             `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomerCreditDetailsFilterType `json:"filters,omitempty" openapi:"$ref:CustomerCreditDetailsFilterType;in:query"`
	Sorts      CustomerCreditDetailsSortType   `json:"sorts,omitempty" openapi:"$ref:CustomerCreditDetailsSortType;in:query"`
}

func (data *CustomersQueryCreditDetailsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"bankAccountNumber": govalidity.Schema{
				"op":    govalidity.New("filter.bankAccountNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.bankAccountNumber.value").Optional(),
			},
			"billingAddressName": govalidity.Schema{
				"op":    govalidity.New("filter.billingAddressName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.billingAddressName.value").Optional(),
			},
		},
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
