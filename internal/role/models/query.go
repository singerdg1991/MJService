package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: RoleSortValue
 */
type RoleSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: RoleSortType
 */
type RoleSortType struct {
	Name      RoleSortValue `json:"name,omitempty" openapi:"$ref:RoleSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Type      RoleSortValue `json:"type,omitempty" openapi:"$ref:RoleSortValue;example:{\"type\":{\"op\":\"asc\"}}"`
	CreatedAt RoleSortValue `json:"created_at,omitempty" openapi:"$ref:RoleSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: RoleFilterType
 */
type RoleFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Type      filters.FilterValue[string] `json:"type,omitempty" openapi:"$ref:FilterValueString;example:{\"type\":{\"op\":\"equals\",\"value\":\"core\"}"`
	CreatedAt filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: RolesQueryRequestParams
 */
type RolesQueryRequestParams struct {
	ID      int            `json:"id,string,omitempty" openapi:"example:1"`
	Page    int            `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters RoleFilterType `json:"filters,omitempty" openapi:"$ref:RoleFilterType;in:query"`
	Sorts   RoleSortType   `json:"sorts,omitempty" openapi:"$ref:RoleSortType;in:query"`
}

func (data *RolesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"type": govalidity.Schema{
				"op":    govalidity.New("filter.type.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.type.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
			"type": govalidity.Schema{
				"op": govalidity.New("sort.type.op"),
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
