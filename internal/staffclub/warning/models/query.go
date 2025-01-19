package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: WarningFilterType
 */
type WarningFilterType struct {
	Title         filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	WarningNumber filters.FilterValue[int]    `json:"warningNumber,omitempty" openapi:"$ref:FilterValueInt;example:{\"warningNumber\":{\"op\":\"equals\",\"value\":1}"`
}

/*
 * @apiDefine: WarningSortValue
 */
type WarningSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: WarningSortType
 */
type WarningSortType struct {
	ID        WarningSortValue `json:"id,omitempty" openapi:"$ref:WarningSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title     WarningSortValue `json:"title,omitempty" openapi:"$ref:WarningSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt WarningSortValue `json:"created_at,omitempty" openapi:"$ref:WarningSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: WarningsQueryRequestParams
 */
type WarningsQueryRequestParams struct {
	ID      int               `json:"id,string,omitempty" openapi:"example:1"`
	Page    int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters WarningFilterType `json:"filters,omitempty" openapi:"$ref:WarningFilterType;in:query"`
	Sorts   WarningSortType   `json:"sorts,omitempty" openapi:"$ref:WarningSortType;in:query"`
}

func (data *WarningsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"warningNumber": govalidity.Schema{
				"op":    govalidity.New("filter.warningNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.warningNumber.value").Optional(),
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
