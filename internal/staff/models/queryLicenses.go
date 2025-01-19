package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: StaffsQueryLicensesFilterType
 */
type StaffsQueryLicensesFilterType struct {
	StaffID   filters.FilterValue[int]    `json:"staffId,omitempty" openapi:"$ref:FilterValueInt;example:{\"staffId\":{\"op\":\"equals\",\"value\":1}}"`
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"name\"}}"`
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"created_at\"}}"`
}

/*
 * @apiDefine: StaffsQueryLicensesSortValue
 */
type StaffsQueryLicensesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: StaffsQueryLicensesSortType
 */
type StaffsQueryLicensesSortType struct {
	Name      StaffsQueryLicensesSortValue `json:"name,omitempty" openapi:"$ref:StaffsQueryLicensesSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	CreatedAt StaffsQueryLicensesSortValue `json:"created_at,omitempty" openapi:"$ref:StaffsQueryLicensesSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: StaffsQueryLicensesRequestParams
 */
type StaffsQueryLicensesRequestParams struct {
	ID      int                           `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                           `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                           `json:"limit,string,omitempty" openapi:"example:10"`
	Filters StaffsQueryLicensesFilterType `json:"filters,omitempty" openapi:"$ref:StaffsQueryLicensesFilterType;in:query"`
	Sorts   StaffsQueryLicensesSortType   `json:"sorts,omitempty" openapi:"$ref:StaffsQueryLicensesSortType;in:query"`
}

func (data *StaffsQueryLicensesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"staffId": govalidity.Schema{
				"op":    govalidity.New("filter.staffId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.staffId.value").Optional(),
			},
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
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
