package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: ServiceTypeFilterType
 */
type ServiceTypeFilterType struct {
	ServiceID   filters.FilterValue[int]    `json:"serviceId,omitempty" openapi:"$ref:FilterValueInt;example:{\"serviceId\":{\"op\":\"equals\",\"value\":1}}"`
	Name        filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: ServiceTypeSortValue
 */
type ServiceTypeSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: ServiceTypeSortType
 */
type ServiceTypeSortType struct {
	Name        ServiceTypeSortValue `json:"name,omitempty" openapi:"$ref:ServiceTypeSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Description ServiceTypeSortValue `json:"description,omitempty" openapi:"$ref:ServiceTypeSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: ServiceTypesQueryRequestParams
 */
type ServiceTypesQueryRequestParams struct {
	ID      int                   `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                   `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                   `json:"limit,string,omitempty" openapi:"example:10"`
	Filters ServiceTypeFilterType `json:"filters,omitempty" openapi:"$ref:ServiceTypeFilterType;in:query"`
	Sorts   ServiceTypeSortType   `json:"sorts,omitempty" openapi:"$ref:ServiceTypeSortType;in:query"`
}

func (data *ServiceTypesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
