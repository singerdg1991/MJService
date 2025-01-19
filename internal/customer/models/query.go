package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CustomerFilterType
 */
type CustomerFilterType struct {
	UserId    filters.FilterValue[int]    `json:"userId,omitempty" openapi:"$ref:FilterValueInt;example:{\"userId\":{\"op\":\"equals\",\"value\":1}"`
	FirstName filters.FilterValue[string] `json:"user.firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"user.firstName\":{\"op\":\"equals\",\"value\":1}"`
	LastName  filters.FilterValue[string] `json:"user.lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"user.lastName\":{\"op\":\"equals\",\"value\":1}"`
}

/*
 * @apiDefine: CustomerSortValue
 */
type CustomerSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomerSortType
 */
type CustomerSortType struct {
	FirstName CustomerSortValue `json:"user.firstName,omitempty" openapi:"$ref:CustomerSortValue;example:{\"user.firstName\":{\"op\":\"asc\"}}"`
	LastName  CustomerSortValue `json:"user.lastName,omitempty" openapi:"$ref:CustomerSortValue;example:{\"user.lastName\":{\"op\":\"asc\"}}"`
	CreatedAt CustomerSortValue `json:"created_at,omitempty" openapi:"$ref:CustomerSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryRequestParams
 */
type CustomersQueryRequestParams struct {
	ID      int                `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int                `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int                `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                `json:"limit,string,omitempty" openapi:"example:10"`
	Filters CustomerFilterType `json:"filters,omitempty" openapi:"$ref:CustomerFilterType;in:query"`
	Sorts   CustomerSortType   `json:"sorts,omitempty" openapi:"$ref:CustomerSortType;in:query"`
}

func (data *CustomersQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"user.firstName": govalidity.Schema{
				"op":    govalidity.New("filter.user.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.firstName.value").Optional(),
			},
			"user.lastName": govalidity.Schema{
				"op":    govalidity.New("filter.user.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.lastName.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"user.firstName": govalidity.Schema{
				"op": govalidity.New("sort.user.firstName.op"),
			},
			"user.lastName": govalidity.Schema{
				"op": govalidity.New("sort.user.lastName.op"),
			},
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
