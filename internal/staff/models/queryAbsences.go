package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: StaffsQueryAbsencesFilterType
 */
type StaffsQueryAbsencesFilterType struct {
	StartDate filters.FilterValue[string] `json:"start_date,omitempty" openapi:"$ref:FilterValueString;example:{\"start_date\":{\"op\":\"equals\",\"value\":\"2020-01-01T00:00:00Z\"}}"`
	EndDate   filters.FilterValue[string] `json:"end_date,omitempty" openapi:"$ref:FilterValueString;example:{\"end_date\":{\"op\":\"equals\",\"value\":\"2020-01-01T00:00:00Z\"}}"`
	Reason    filters.FilterValue[string] `json:"reason,omitempty" openapi:"$ref:FilterValueString;example:{\"reason\":{\"op\":\"equals\",\"value\":\"reason\"}}"`
	Status    filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"pending\"}}"`
}

/*
 * @apiDefine: StaffsQueryAbsencesSortValue
 */
type StaffsQueryAbsencesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: StaffsQueryAbsencesSortType
 */
type StaffsQueryAbsencesSortType struct {
	StartDate StaffsQueryAbsencesSortValue `json:"start_date,omitempty" openapi:"$ref:StaffsQueryAbsencesSortValue;example:{\"start_date\":{\"op\":\"asc\"}}"`
	EndDate   StaffsQueryAbsencesSortValue `json:"end_date,omitempty" openapi:"$ref:StaffsQueryAbsencesSortValue;example:{\"end_date\":{\"op\":\"asc\"}}"`
	Reason    StaffsQueryAbsencesSortValue `json:"reason,omitempty" openapi:"$ref:StaffsQueryAbsencesSortValue;example:{\"reason\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: StaffsQueryAbsencesRequestParams
 */
type StaffsQueryAbsencesRequestParams struct {
	ID      int                           `json:"id,string,omitempty" openapi:"example:1"`
	StaffID int                           `json:"staffId,string,omitempty" openapi:"example:1;required"`
	Page    int                           `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                           `json:"limit,string,omitempty" openapi:"example:10"`
	Filters StaffsQueryAbsencesFilterType `json:"filters,omitempty" openapi:"$ref:StaffsQueryAbsencesFilterType;in:query"`
	Sorts   StaffsQueryAbsencesSortType   `json:"sorts,omitempty" openapi:"$ref:StaffsQueryAbsencesSortType;in:query"`
}

func (data *StaffsQueryAbsencesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"staffId": govalidity.New("staffId").Int().Required(),
		"page":    govalidity.New("page").Int().Default("1"),
		"limit":   govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"start_date": govalidity.Schema{
				"op":    govalidity.New("filter.start_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_date.value").Optional(),
			},
			"end_date": govalidity.Schema{
				"op":    govalidity.New("filter.end_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.end_date.value").Optional(),
			},
			"reason": govalidity.Schema{
				"op":    govalidity.New("filter.reason.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.reason.value").Optional(),
			},
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
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
