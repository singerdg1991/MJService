package models

import (
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysRequestParamsFilterType
 */
type CyclesQueryShiftCustomerHomeKeysRequestParamsFilterType struct {
	KeyNo  filters.FilterValue[string] `json:"keyNo,omitempty" openapi:"$ref:FilterValueString;example:{\"keyNo\":{\"op\":\"equals\",\"value\":\"123456\"}"`
	Status filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"accepted\"}"`
}

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysRequestParamsSortValue
 */
type CyclesQueryShiftCustomerHomeKeysRequestParamsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysRequestParamsSortType
 */
type CyclesQueryShiftCustomerHomeKeysRequestParamsSortType struct {
	CreatedAt CyclesQueryShiftCustomerHomeKeysRequestParamsSortValue `json:"created_at,omitempty" openapi:"$ref:CyclesQueryShiftCustomerHomeKeysRequestParamsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysRequestParams
 */
type CyclesQueryShiftCustomerHomeKeysRequestParams struct {
	ID      int                                                     `json:"id,string,omitempty" openapi:"example:1"`
	ShiftID int                                                     `json:"shiftid,string,omitempty" openapi:"example:1"`
	Page    int                                                     `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                                                     `json:"limit,string,omitempty" openapi:"example:10"`
	Filters CyclesQueryShiftCustomerHomeKeysRequestParamsFilterType `json:"filters,omitempty" openapi:"$ref:CyclesQueryShiftCustomerHomeKeysRequestParamsFilterType;in:query"`
	Sorts   CyclesQueryShiftCustomerHomeKeysRequestParamsSortType   `json:"sorts,omitempty" openapi:"$ref:CyclesQueryShiftCustomerHomeKeysRequestParamsSortType;in:query"`
}

// ValidateQueries validates the CyclesQueryShiftCustomerHomeKeysRequestParams based on the provided schema and request.
//
// It takes an http.Request as a parameter to validate the query parameters.
// It returns a govalidity.ValidityResponseErrors if the validation fails, otherwise nil.
func (data *CyclesQueryShiftCustomerHomeKeysRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"shiftid": govalidity.New("shiftid").Int().Optional(),
		"page":    govalidity.New("page").Int().Default("1"),
		"limit":   govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
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
