package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: TodoFilterType
 */
type TodoFilterType struct {
	Title       filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Date        filters.FilterValue[string] `json:"date,omitempty" openapi:"$ref:FilterValueString;example:{\"date\":{\"op\":\"equals\",\"value\":\"2001-02-06\"}"`
	Time        filters.FilterValue[string] `json:"time,omitempty" openapi:"$ref:FilterValueString;example:{\"time\":{\"op\":\"equals\",\"value\":\"12:00\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Status      filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"done\"}"`
}

/*
 * @apiDefine: TodoSortValue
 */
type TodoSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: TodoSortType
 */
type TodoSortType struct {
	Title       TodoSortValue `json:"title,omitempty" openapi:"$ref:TodoSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	Date        TodoSortValue `json:"date,omitempty" openapi:"$ref:TodoSortValue;example:{\"date\":{\"op\":\"asc\"}}"`
	Time        TodoSortValue `json:"time,omitempty" openapi:"$ref:TodoSortValue;example:{\"time\":{\"op\":\"asc\"}}"`
	UserID      TodoSortValue `json:"userId,omitempty" openapi:"$ref:TodoSortValue;example:{\"userId\":{\"op\":\"asc\"}}"`
	Description TodoSortValue `json:"description,omitempty" openapi:"$ref:TodoSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	CreatedAt   TodoSortValue `json:"created_at,omitempty" openapi:"$ref:TodoSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: TodosQueryRequestParams
 */
type TodosQueryRequestParams struct {
	ID      int            `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int            `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int            `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters TodoFilterType `json:"filters,omitempty" openapi:"$ref:TodoFilterType;in:query"`
	Sorts   TodoSortType   `json:"sorts,omitempty" openapi:"$ref:TodoSortType;in:query"`
}

func (data *TodosQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"date": govalidity.Schema{
				"op":    govalidity.New("filter.date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.date.value").Optional(),
			},
			"time": govalidity.Schema{
				"op":    govalidity.New("filter.time.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.time.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
			},
			"date": govalidity.Schema{
				"op": govalidity.New("sort.date.op"),
			},
			"time": govalidity.Schema{
				"op": govalidity.New("sort.time.op"),
			},
			"userId": govalidity.Schema{
				"op": govalidity.New("sort.userId.op"),
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
