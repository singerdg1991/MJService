package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CustomersQueryAbsencesFilterType
 */
type CustomersQueryAbsencesFilterType struct {
	StartDate filters.FilterValue[string] `json:"start_date,omitempty" openapi:"$ref:FilterValueString;example:{\"start_date\":{\"op\":\"equals\",\"value\":\"2020-01-01T00:00:00Z\"}}"`
	EndDate   filters.FilterValue[string] `json:"end_date,omitempty" openapi:"$ref:FilterValueString;example:{\"end_date\":{\"op\":\"equals\",\"value\":\"2020-01-01T00:00:00Z\"}}"`
	Reason    filters.FilterValue[string] `json:"reason,omitempty" openapi:"$ref:FilterValueString;example:{\"reason\":{\"op\":\"equals\",\"value\":\"reason\"}}"`
}

/*
 * @apiDefine: CustomersQueryAbsencesSortValue
 */
type CustomersQueryAbsencesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryAbsencesSortType
 */
type CustomersQueryAbsencesSortType struct {
	StartDate CustomersQueryAbsencesSortValue `json:"start_date,omitempty" openapi:"$ref:CustomersQueryAbsencesSortValue;example:{\"start_date\":{\"op\":\"asc\"}}"`
	EndDate   CustomersQueryAbsencesSortValue `json:"end_date,omitempty" openapi:"$ref:CustomersQueryAbsencesSortValue;example:{\"end_date\":{\"op\":\"asc\"}}"`
	Reason    CustomersQueryAbsencesSortValue `json:"reason,omitempty" openapi:"$ref:CustomersQueryAbsencesSortValue;example:{\"reason\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryAbsencesRequestParams
 */
type CustomersQueryAbsencesRequestParams struct {
	ID         int                              `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                              `json:"customerId,string,omitempty" openapi:"example:1;required"`
	Page       int                              `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                              `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryAbsencesFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryAbsencesFilterType;in:query"`
	Sorts      CustomersQueryAbsencesSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryAbsencesSortType;in:query"`
}

func (data *CustomersQueryAbsencesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Required(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"start_date": govalidity.Schema{
				"op":    govalidity.New("filter.start_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_date.value").Optional(),
			},
			"end_date": govalidity.Schema{
				"op":    govalidity.New("filter.end_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.end_date.value").Optional(),
			},
			"reason": govalidity.Schema{
				"op":    govalidity.New("filter.reason.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.reason.value").Optional(),
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
