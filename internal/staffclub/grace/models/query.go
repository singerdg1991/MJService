package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: GraceFilterType
 */
type GraceFilterType struct {
	Title       filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	GraceNumber filters.FilterValue[int]    `json:"graceNumber,omitempty" openapi:"$ref:FilterValueInt;example:{\"graceNumber\":{\"op\":\"equals\",\"value\":1}"`
}

/*
 * @apiDefine: GraceSortValue
 */
type GraceSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: GraceSortType
 */
type GraceSortType struct {
	ID        GraceSortValue `json:"id,omitempty" openapi:"$ref:GraceSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title     GraceSortValue `json:"title,omitempty" openapi:"$ref:GraceSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt GraceSortValue `json:"created_at,omitempty" openapi:"$ref:GraceSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: GracesQueryRequestParams
 */
type GracesQueryRequestParams struct {
	ID      int             `json:"id,string,omitempty" openapi:"example:1"`
	Page    int             `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters GraceFilterType `json:"filters,omitempty" openapi:"$ref:GraceFilterType;in:query"`
	Sorts   GraceSortType   `json:"sorts,omitempty" openapi:"$ref:GraceSortType;in:query"`
}

func (data *GracesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"graceNumber": govalidity.Schema{
				"op":    govalidity.New("filter.graceNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.graceNumber.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
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
