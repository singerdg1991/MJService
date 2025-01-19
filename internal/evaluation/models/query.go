package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: EvaluationFilterType
 */
type EvaluationFilterType struct {
	StaffID        filters.FilterValue[int]    `json:"staffId,omitempty" openapi:"$ref:FilterValueInt;example:{\"staffId\":{\"op\":\"equals\",\"value\":1}}"`
	EvaluationType filters.FilterValue[string] `json:"evaluationType,omitempty" openapi:"$ref:FilterValueString;example:{\"evaluationType\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Title          filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description    filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: EvaluationSortValue
 */
type EvaluationSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: EvaluationSortType
 */
type EvaluationSortType struct {
	EvaluationType EvaluationSortValue `json:"evaluationType,omitempty" openapi:"$ref:EvaluationSortValue;example:{\"evaluationType\":{\"op\":\"asc\"}}"`
	Title          EvaluationSortValue `json:"title,omitempty" openapi:"$ref:EvaluationSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	Description    EvaluationSortValue `json:"description,omitempty" openapi:"$ref:EvaluationSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	CreatedAt      EvaluationSortValue `json:"created_at,omitempty" openapi:"$ref:EvaluationSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: EvaluationsQueryRequestParams
 */
type EvaluationsQueryRequestParams struct {
	ID      int                  `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                  `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                  `json:"limit,string,omitempty" openapi:"example:10"`
	Filters EvaluationFilterType `json:"filters,omitempty" openapi:"$ref:EvaluationFilterType;in:query"`
	Sorts   EvaluationSortType   `json:"sorts,omitempty" openapi:"$ref:EvaluationSortType;in:query"`
}

func (data *EvaluationsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"staffId": govalidity.Schema{
				"op":    govalidity.New("filter.staffId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.staffId.value").Optional(),
			},
			"evaluationType": govalidity.Schema{
				"op":    govalidity.New("filter.evaluationType.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.evaluationType.value").Optional(),
			},
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"evaluationType": govalidity.Schema{
				"op": govalidity.New("sort.evaluationType.op"),
			},
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
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
