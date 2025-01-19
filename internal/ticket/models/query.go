package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: TicketFilterType
 */
type TicketFilterType struct {
	Title filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"title\"}"`
}

/*
 * @apiDefine: TicketSortValue
 */
type TicketSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: TicketSortType
 */
type TicketSortType struct {
	ID    TicketSortValue `json:"id,omitempty" openapi:"$ref:TicketSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title TicketSortValue `json:"title,omitempty" openapi:"$ref:TicketSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: TicketsQueryRequestParams
 */
type TicketsQueryRequestParams struct {
	ID      int              `json:"id,string,omitempty" openapi:"example:1"`
	Page    int              `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int              `json:"limit,string,omitempty" openapi:"example:10"`
	Filters TicketFilterType `json:"filters,omitempty" openapi:"$ref:TicketFilterType;in:query"`
	Sorts   TicketSortType   `json:"sorts,omitempty" openapi:"$ref:TicketSortType;in:query"`
}

func (data *TicketsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
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
