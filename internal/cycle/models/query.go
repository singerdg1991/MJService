package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CycleFilterType
 */
type CycleFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	StartDate filters.FilterValue[string] `json:"start_date,omitempty" openapi:"$ref:FilterValueString;example:{\"start_date\":{\"op\":\"equals\",\"value\":\"2024-08-10\"}"`
}

/*
 * @apiDefine: CycleSortValue
 */
type CycleSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CycleSortType
 */
type CycleSortType struct {
	ID        CycleSortValue `json:"id,omitempty" openapi:"$ref:CycleSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Name      CycleSortValue `json:"name,omitempty" openapi:"$ref:CycleSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	CreatedAt CycleSortValue `json:"created_at,omitempty" openapi:"$ref:CycleSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CyclesQueryRequestParams
 */
type CyclesQueryRequestParams struct {
	ID        int             `json:"id,string,omitempty" openapi:"example:1"`
	SectionID int             `json:"sectionid,string,omitempty" openapi:"example:1"`
	Page      int             `json:"page,string,omitempty" openapi:"example:1"`
	Limit     int             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters   CycleFilterType `json:"filters,omitempty" openapi:"$ref:CycleFilterType;in:query"`
	Sorts     CycleSortType   `json:"sorts,omitempty" openapi:"$ref:CycleSortType;in:query"`
}

// ValidateQueries validates the CyclesQueryRequestParams based on the provided schema and request.
//
// It takes an http.Request as a parameter to validate the query parameters.
// It returns a govalidity.ValidityResponseErrors if the validation fails, otherwise nil.
func (data *CyclesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":        govalidity.New("id").Int().Optional(),
		"sectionid": govalidity.New("sectionid").Int().Optional(),
		"page":      govalidity.New("page").Int().Default("1"),
		"limit":     govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"start_date": govalidity.Schema{
				"op":    govalidity.New("filter.start_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_date.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Transform start_date
	if data.Filters.StartDate.Value != "" {
		// Convert to 2021-08-10T00:00:00Z format
		startDate, err := time.Parse("2006-01-02", data.Filters.StartDate.Value)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"start_date": []string{"start_date must be a valid date in format YYYY-MM-DD"},
			}
		}
		data.Filters.StartDate.Value = startDate.Format(time.RFC3339)
	}

	return nil
}
