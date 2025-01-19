package models

import (
	"net/http"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CycleQueryCycleShiftsFilterType
 */
type CycleQueryCycleShiftsFilterType struct {
	ShiftName filters.FilterValue[string] `json:"shiftName,omitempty" openapi:"$ref:FilterValueString;example:{\"shiftName\":{\"op\":\"equals\",\"value\":\"morning\"}"`
	DateTime  filters.FilterValue[string] `json:"dateTime,omitempty" openapi:"$ref:FilterValueString;example:{\"dateTime\":{\"op\":\"equals\",\"value\":\"2020-01-01\"}"`
	Status    filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"not-started\"}"`
}

/*
 * @apiDefine: CycleQueryCycleShiftsSortValue
 */
type CycleQueryCycleShiftsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CycleQueryCycleShiftsSortType
 */
type CycleQueryCycleShiftsSortType struct {
}

/*
 * @apiDefine: CyclesQueryCycleShiftsRequestParams
 */
type CyclesQueryCycleShiftsRequestParams struct {
	ID      int                             `json:"id,string,omitempty" openapi:"example:1"`
	CycleID int                             `json:"cycleid,string,omitempty" openapi:"example:1"`
	Page    int                             `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters CycleQueryCycleShiftsFilterType `json:"filters,omitempty" openapi:"$ref:CycleQueryCycleShiftsFilterType;in:query"`
	Sorts   CycleQueryCycleShiftsSortType   `json:"sorts,omitempty" openapi:"$ref:CycleQueryCycleShiftsSortType;in:query"`
}

// ValidateQueries validates the queries of the CyclesQueryCycleShiftsRequestParams.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesQueryCycleShiftsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"cycleid": govalidity.New("cycleid").Int().Optional(),
		"page":    govalidity.New("page").Int().Default("1"),
		"limit":   govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"shiftName": govalidity.Schema{
				"op":    govalidity.New("filter.shiftName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.shiftName.value").Optional(),
			},
			"datetime": govalidity.Schema{
				"op":    govalidity.New("filter.datetime.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.datetime.value").Optional(),
			},
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	datetime := data.Filters.DateTime.Value
	if datetime != "" {
		_, err := time.Parse("2006-01-02", datetime)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"filter.datetime.value": []string{"datetime should have format YYYY-MM-DD"},
			}
		}
	}

	return nil
}
