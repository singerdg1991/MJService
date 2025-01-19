package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: DiagnoseFilterType
 */
type DiagnoseFilterType struct {
	Title       filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Code        filters.FilterValue[string] `json:"code,omitempty" openapi:"$ref:FilterValueString;example:{\"code\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt   filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: DiagnoseSortValue
 */
type DiagnoseSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: DiagnoseSortType
 */
type DiagnoseSortType struct {
	Title       DiagnoseSortValue `json:"title,omitempty" openapi:"$ref:DiagnoseSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	Code        DiagnoseSortValue `json:"code,omitempty" openapi:"$ref:DiagnoseSortValue;example:{\"code\":{\"op\":\"asc\"}}"`
	Description DiagnoseSortValue `json:"description,omitempty" openapi:"$ref:DiagnoseSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	CreatedAt   DiagnoseSortValue `json:"created_at,omitempty" openapi:"$ref:DiagnoseSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: DiagnosesQueryRequestParams
 */
type DiagnosesQueryRequestParams struct {
	ID      int                `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                `json:"limit,string,omitempty" openapi:"example:10"`
	Filters DiagnoseFilterType `json:"filters,omitempty" openapi:"$ref:DiagnoseFilterType;in:query"`
	Sorts   DiagnoseSortType   `json:"sorts,omitempty" openapi:"$ref:DiagnoseSortType;in:query"`
}

func (data *DiagnosesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"code": govalidity.Schema{
				"op":    govalidity.New("filter.code.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.code.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
			},
			"code": govalidity.Schema{
				"op": govalidity.New("sort.code.op"),
			},
			"description": govalidity.Schema{
				"op": govalidity.New("sort.description.op"),
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
