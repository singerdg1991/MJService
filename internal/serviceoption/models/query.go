package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: ServiceOptionFilterType
 */
type ServiceOptionFilterType struct {
	ServiceTypeID filters.FilterValue[int]    `json:"serviceTypeId,omitempty" openapi:"$ref:FilterValueInt;example:{\"serviceTypeId\":{\"op\":\"equals\",\"value\":1}}"`
	Name          filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description   filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: ServiceOptionSortValue
 */
type ServiceOptionSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: ServiceOptionSortType
 */
type ServiceOptionSortType struct {
	Name        ServiceOptionSortValue `json:"name,omitempty" openapi:"$ref:ServiceOptionSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Description ServiceOptionSortValue `json:"description,omitempty" openapi:"$ref:ServiceOptionSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: ServiceOptionsQueryRequestParams
 */
type ServiceOptionsQueryRequestParams struct {
	ID      int                     `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                     `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                     `json:"limit,string,omitempty" openapi:"example:10"`
	Filters ServiceOptionFilterType `json:"filters,omitempty" openapi:"$ref:ServiceOptionFilterType;in:query"`
	Sorts   ServiceOptionSortType   `json:"sorts,omitempty" openapi:"$ref:ServiceOptionSortType;in:query"`
}

func (data *ServiceOptionsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
