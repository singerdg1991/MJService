package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: PermissionFilterType
 */
type PermissionFilterType struct {
	Name  filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Title filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: PermissionSortValue
 */
type PermissionSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: PermissionSortType
 */
type PermissionSortType struct {
	ID        PermissionSortValue `json:"id,omitempty" openapi:"$ref:PermissionSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Name      PermissionSortValue `json:"name,omitempty" openapi:"$ref:PermissionSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Title     PermissionSortValue `json:"title,omitempty" openapi:"$ref:PermissionSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt PermissionSortValue `json:"created_at,omitempty" openapi:"$ref:PermissionSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: PermissionsQueryRequestParams
 */
type PermissionsQueryRequestParams struct {
	ID      int                  `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                  `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                  `json:"limit,string,omitempty" openapi:"example:10"`
	Filters PermissionFilterType `json:"filters,omitempty" openapi:"$ref:PermissionFilterType;in:query"`
	Sorts   PermissionSortType   `json:"sorts,omitempty" openapi:"$ref:PermissionSortType;in:query"`
}

func (data *PermissionsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
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
