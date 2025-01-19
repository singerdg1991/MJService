package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: PunishmentFilterType
 */
type PunishmentFilterType struct {
	Name        filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: PunishmentSortValue
 */
type PunishmentSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: PunishmentSortType
 */
type PunishmentSortType struct {
	Name        PunishmentSortValue `json:"name,omitempty" openapi:"$ref:PunishmentSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Description PunishmentSortValue `json:"description,omitempty" openapi:"$ref:PunishmentSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: PunishmentsQueryRequestParams
 */
type PunishmentsQueryRequestParams struct {
	ID      int                  `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                  `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                  `json:"limit,string,omitempty" openapi:"example:10"`
	Filters PunishmentFilterType `json:"filters,omitempty" openapi:"$ref:PunishmentFilterType;in:query"`
	Sorts   PunishmentSortType   `json:"sorts,omitempty" openapi:"$ref:PunishmentSortType;in:query"`
}

func (data *PunishmentsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
			"description": govalidity.Schema{
				"op": govalidity.New("sort.description.op"),
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
