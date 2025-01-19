package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: EquipmentFilterType
 */
type EquipmentFilterType struct {
	Name        filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: EquipmentSortValue
 */
type EquipmentSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: EquipmentSortType
 */
type EquipmentSortType struct {
	Name        EquipmentSortValue `json:"name,omitempty" openapi:"$ref:EquipmentSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Description EquipmentSortValue `json:"description,omitempty" openapi:"$ref:EquipmentSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: EquipmentsQueryRequestParams
 */
type EquipmentsQueryRequestParams struct {
	ID      int                 `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                 `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                 `json:"limit,string,omitempty" openapi:"example:10"`
	Filters EquipmentFilterType `json:"filters,omitempty" openapi:"$ref:EquipmentFilterType;in:query"`
	Sorts   EquipmentSortType   `json:"sorts,omitempty" openapi:"$ref:EquipmentSortType;in:query"`
}

func (data *EquipmentsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
