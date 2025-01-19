package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: TrashFilterType
 */
type TrashFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: TrashesQueryRequestParams
 */
type TrashesQueryRequestParams struct {
	ID      int             `json:"id,string,omitempty" openapi:"example:1"`
	Page    int             `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters TrashFilterType `json:"filters,omitempty" openapi:"$ref:TrashFilterType;in:query"`
}

func (data *TrashesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
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
