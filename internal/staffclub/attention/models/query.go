package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: AttentionFilterType
 */
type AttentionFilterType struct {
	Title           filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	AttentionNumber filters.FilterValue[int]    `json:"attentionNumber,omitempty" openapi:"$ref:FilterValueInt;example:{\"attentionNumber\":{\"op\":\"equals\",\"value\":1}"`
}

/*
 * @apiDefine: AttentionSortValue
 */
type AttentionSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: AttentionSortType
 */
type AttentionSortType struct {
	ID        AttentionSortValue `json:"id,omitempty" openapi:"$ref:AttentionSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title     AttentionSortValue `json:"title,omitempty" openapi:"$ref:AttentionSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt AttentionSortValue `json:"created_at,omitempty" openapi:"$ref:AttentionSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: AttentionsQueryRequestParams
 */
type AttentionsQueryRequestParams struct {
	ID      int                 `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                 `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                 `json:"limit,string,omitempty" openapi:"example:10"`
	Filters AttentionFilterType `json:"filters,omitempty" openapi:"$ref:AttentionFilterType;in:query"`
	Sorts   AttentionSortType   `json:"sorts,omitempty" openapi:"$ref:AttentionSortType;in:query"`
}

func (data *AttentionsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"attentionNumber": govalidity.Schema{
				"op":    govalidity.New("filter.attentionNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.attentionNumber.value").Optional(),
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
