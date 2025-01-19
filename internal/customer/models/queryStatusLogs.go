package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CustomersQueryStatusLogsFilterType
 */
type CustomersQueryStatusLogsFilterType struct {
	StatusValue filters.FilterValue[string] `json:"statusValue,omitempty" openapi:"$ref:FilterValueString;example:{\"statusValue\":{\"op\":\"equals\",\"value\":\"statusValue\"}}"`
	FirstName   filters.FilterValue[string] `json:"firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"firstName\":{\"op\":\"equals\",\"value\":\"firstName\"}}"`
	LastName    filters.FilterValue[string] `json:"lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"lastName\":{\"op\":\"equals\",\"value\":\"lastName\"}}"`
	FullName    filters.FilterValue[string] `json:"fullName,omitempty" openapi:"$ref:FilterValueString;example:{\"fullName\":{\"op\":\"equals\",\"value\":\"fullName\"}}"`
	CreatedAt   filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"2021-01-01T00:00:00Z\"}}"`
}

/*
 * @apiDefine: CustomersQueryStatusLogsSortValue
 */
type CustomersQueryStatusLogsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryStatusLogsSortType
 */
type CustomersQueryStatusLogsSortType struct {
	ID          CustomersQueryStatusLogsSortValue `json:"id,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	StatusValue CustomersQueryStatusLogsSortValue `json:"statusValue,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortValue;example:{\"statusValue\":{\"op\":\"asc\"}}"`
	FirstName   CustomersQueryStatusLogsSortValue `json:"firstName,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortValue;example:{\"firstName\":{\"op\":\"asc\"}}"`
	LastName    CustomersQueryStatusLogsSortValue `json:"lastName,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortValue;example:{\"lastName\":{\"op\":\"asc\"}}"`
	CreatedAt   CustomersQueryStatusLogsSortValue `json:"created_at,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryStatusLogsRequestParams
 */
type CustomersQueryStatusLogsRequestParams struct {
	ID         int                                `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                                `json:"customerId,string,omitempty" openapi:"example:1;required"`
	Page       int                                `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                                `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryStatusLogsFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryStatusLogsFilterType;in:query"`
	Sorts      CustomersQueryStatusLogsSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryStatusLogsSortType;in:query"`
}

func (data *CustomersQueryStatusLogsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Required(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"beforeValue": govalidity.Schema{
				"op":    govalidity.New("filter.beforeValue.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.beforeValue.value").Optional(),
			},
			"statusValue": govalidity.Schema{
				"op":    govalidity.New("filter.statusValue.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.statusValue.value").Optional(),
			},
			"firstName": govalidity.Schema{
				"op":    govalidity.New("filter.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.firstName.value").Optional(),
			},
			"lastName": govalidity.Schema{
				"op":    govalidity.New("filter.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.lastName.value").Optional(),
			},
			"fullName": govalidity.Schema{
				"op":    govalidity.New("filter.fullName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.fullName.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"beforeValue": govalidity.Schema{
				"op": govalidity.New("sort.beforeValue.op"),
			},
			"statusValue": govalidity.Schema{
				"op": govalidity.New("sort.statusValue.op"),
			},
			"firstName": govalidity.Schema{
				"op": govalidity.New("sort.firstName.op"),
			},
			"lastName": govalidity.Schema{
				"op": govalidity.New("sort.lastName.op"),
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
