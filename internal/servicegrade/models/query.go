package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: ServiceGradeFilterType
 */
type ServiceGradeFilterType struct {
	Name        filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Grade       filters.FilterValue[int]    `json:"grade,omitempty" openapi:"$ref:FilterValueInt;example:{\"grade\":{\"op\":\"equals\",\"value\":1}"`
	Color       filters.FilterValue[string] `json:"color,omitempty" openapi:"$ref:FilterValueString;example:{\"color\":{\"op\":\"equals\",\"value\":\"#000000\"}"`
}

/*
 * @apiDefine: ServiceGradeSortValue
 */
type ServiceGradeSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: ServiceGradeSortType
 */
type ServiceGradeSortType struct {
	Name        ServiceGradeSortValue `json:"name,omitempty" openapi:"$ref:ServiceGradeSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Description ServiceGradeSortValue `json:"description,omitempty" openapi:"$ref:ServiceGradeSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	Grade       ServiceGradeSortValue `json:"grade,omitempty" openapi:"$ref:ServiceGradeSortValue;example:{\"grade\":{\"op\":\"asc\"}}"`
	Color       ServiceGradeSortValue `json:"color,omitempty" openapi:"$ref:ServiceGradeSortValue;example:{\"color\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: ServiceGradesQueryRequestParams
 */
type ServiceGradesQueryRequestParams struct {
	ID      int                    `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                    `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                    `json:"limit,string,omitempty" openapi:"example:10"`
	Filters ServiceGradeFilterType `json:"filters,omitempty" openapi:"$ref:ServiceGradeFilterType;in:query"`
	Sorts   ServiceGradeSortType   `json:"sorts,omitempty" openapi:"$ref:ServiceGradeSortType;in:query"`
}

func (data *ServiceGradesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"grade": govalidity.Schema{
				"op":    govalidity.New("filter.grade.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.grade.value").Optional(),
			},
			"color": govalidity.Schema{
				"op":    govalidity.New("filter.color.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.color.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
			"description": govalidity.Schema{
				"op": govalidity.New("sort.description.op"),
			},
			"grade": govalidity.Schema{
				"op": govalidity.New("sort.grade.op"),
			},
			"color": govalidity.Schema{
				"op": govalidity.New("sort.color.op"),
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
