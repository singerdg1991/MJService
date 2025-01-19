package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: MedicineFilterType
 */
type MedicineFilterType struct {
	Name         filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Code         filters.FilterValue[string] `json:"code,omitempty" openapi:"$ref:FilterValueString;example:{\"code\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Availability filters.FilterValue[string] `json:"availability,omitempty" openapi:"$ref:FilterValueString;example:{\"availability\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Manufacturer filters.FilterValue[string] `json:"manufacturer,omitempty" openapi:"$ref:FilterValueString;example:{\"manufacturer\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	PurposeOfUse filters.FilterValue[string] `json:"purposeOfUse,omitempty" openapi:"$ref:FilterValueString;example:{\"purposeOfUse\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Instruction  filters.FilterValue[string] `json:"instruction,omitempty" openapi:"$ref:FilterValueString;example:{\"instruction\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	SideEffects  filters.FilterValue[string] `json:"sideEffects,omitempty" openapi:"$ref:FilterValueString;example:{\"sideEffects\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Conditions   filters.FilterValue[string] `json:"conditions,omitempty" openapi:"$ref:FilterValueString;example:{\"conditions\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description  filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: MedicineSortValue
 */
type MedicineSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: MedicineSortType
 */
type MedicineSortType struct {
	ID           MedicineSortValue `json:"id,omitempty" openapi:"$ref:MedicineSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Name         MedicineSortValue `json:"name,omitempty" openapi:"$ref:MedicineSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Code         MedicineSortValue `json:"code,omitempty" openapi:"$ref:MedicineSortValue;example:{\"code\":{\"op\":\"asc\"}}"`
	Availability MedicineSortValue `json:"availability,omitempty" openapi:"$ref:MedicineSortValue;example:{\"availability\":{\"op\":\"asc\"}}"`
	Manufacturer MedicineSortValue `json:"manufacturer,omitempty" openapi:"$ref:MedicineSortValue;example:{\"manufacturer\":{\"op\":\"asc\"}}"`
	PurposeOfUse MedicineSortValue `json:"purposeOfUse,omitempty" openapi:"$ref:MedicineSortValue;example:{\"purposeOfUse\":{\"op\":\"asc\"}}"`
	Instruction  MedicineSortValue `json:"instruction,omitempty" openapi:"$ref:MedicineSortValue;example:{\"instruction\":{\"op\":\"asc\"}}"`
	SideEffects  MedicineSortValue `json:"sideEffects,omitempty" openapi:"$ref:MedicineSortValue;example:{\"sideEffects\":{\"op\":\"asc\"}}"`
	Conditions   MedicineSortValue `json:"conditions,omitempty" openapi:"$ref:MedicineSortValue;example:{\"conditions\":{\"op\":\"asc\"}}"`
	Description  MedicineSortValue `json:"description,omitempty" openapi:"$ref:MedicineSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
	CreatedAt    MedicineSortValue `json:"created_at,omitempty" openapi:"$ref:MedicineSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: MedicinesQueryRequestParams
 */
type MedicinesQueryRequestParams struct {
	ID      int                `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                `json:"limit,string,omitempty" openapi:"example:10"`
	Filters MedicineFilterType `json:"filters,omitempty" openapi:"$ref:MedicineFilterType;in:query"`
	Sorts   MedicineSortType   `json:"sorts,omitempty" openapi:"$ref:MedicineSortType;in:query"`
}

func (data *MedicinesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"code": govalidity.Schema{
				"op":    govalidity.New("filter.code.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.code.value").Optional(),
			},
			"availability": govalidity.Schema{
				"op":    govalidity.New("filter.availability.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.availability.value").Optional(),
			},
			"manufacturer": govalidity.Schema{
				"op":    govalidity.New("filter.manufacturer.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.manufacturer.value").Optional(),
			},
			"purposeOfUse": govalidity.Schema{
				"op":    govalidity.New("filter.purposeOfUse.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.purposeOfUse.value").Optional(),
			},
			"instruction": govalidity.Schema{
				"op":    govalidity.New("filter.instruction.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.instruction.value").Optional(),
			},
			"sideEffects": govalidity.Schema{
				"op":    govalidity.New("filter.sideEffects.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.sideEffects.value").Optional(),
			},
			"conditions": govalidity.Schema{
				"op":    govalidity.New("filter.conditions.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.conditions.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
			"code": govalidity.Schema{
				"op": govalidity.New("sort.code.op"),
			},
			"availability": govalidity.Schema{
				"op": govalidity.New("sort.availability.op"),
			},
			"manufacturer": govalidity.Schema{
				"op": govalidity.New("sort.manufacturer.op"),
			},
			"purposeOfUse": govalidity.Schema{
				"op": govalidity.New("sort.purposeOfUse.op"),
			},
			"instruction": govalidity.Schema{
				"op": govalidity.New("sort.instruction.op"),
			},
			"sideEffects": govalidity.Schema{
				"op": govalidity.New("sort.sideEffects.op"),
			},
			"conditions": govalidity.Schema{
				"op": govalidity.New("sort.conditions.op"),
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
