package models

import (
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CyclesQueryVisitsTodosFilterType
 */
type CyclesQueryVisitsTodosFilterType struct {
	Status filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"done\"}"`
}

/*
 * @apiDefine: CyclesQueryVisitsTodosSortValue
 */
type CyclesQueryVisitsTodosSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CyclesQueryVisitsTodosSortType
 */
type CyclesQueryVisitsTodosSortType struct {
	CreatedAt CyclesQueryVisitsTodosSortValue `json:"created_at,omitempty" openapi:"$ref:CyclesQueryVisitsTodosSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CyclesQueryVisitsTodosRequestParams
 */
type CyclesQueryVisitsTodosRequestParams struct {
	ID                 int                              `json:"id,string,omitempty" openapi:"example:1"`
	CyclePickupShiftID int                              `json:"cyclepickupshiftid,string,omitempty" openapi:"example:1"`
	Page               int                              `json:"page,string,omitempty" openapi:"example:1"`
	Limit              int                              `json:"limit,string,omitempty" openapi:"example:10"`
	Filters            CyclesQueryVisitsTodosFilterType `json:"filters,omitempty" openapi:"$ref:CyclesQueryVisitsTodosFilterType;in:query"`
	Sorts              CyclesQueryVisitsTodosSortType   `json:"sorts,omitempty" openapi:"$ref:CyclesQueryVisitsTodosSortType;in:query"`
}

// ValidateQueries validates the CyclesQueryVisitsTodosRequestParams based on the provided schema and request.
//
// It takes an http.Request as a parameter to validate the query parameters.
// It returns a govalidity.ValidityResponseErrors if the validation fails, otherwise nil.
func (data *CyclesQueryVisitsTodosRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":                 govalidity.New("id").Int().Optional(),
		"cyclepickupshiftid": govalidity.New("cyclepickupshiftid").Int().Optional(),
		"page":               govalidity.New("page").Int().Default("1"),
		"limit":              govalidity.New("limit").Int().Default("10"),
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

	// Check if the request queries has error
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
