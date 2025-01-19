package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: KeikkalasQueryShiftStatisticsFilterType
 */
type KeikkalasQueryShiftStatisticsFilterType struct {
	StartDate filters.FilterValue[string] `json:"start_date,omitempty" openapi:"$ref:FilterValueString;example:{\"start_date\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
}

/*
 * @apiDefine: KeikkalasQueryShiftStatisticsRequestParams
 */
type KeikkalasQueryShiftStatisticsRequestParams struct {
	Filters KeikkalasQueryShiftStatisticsFilterType `json:"filters,omitempty" openapi:"$ref:KeikkalasQueryShiftStatisticsFilterType;in:query"`
}

func (data *KeikkalasQueryShiftStatisticsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"filters": govalidity.Schema{
			"start_date": govalidity.Schema{
				"op":    govalidity.New("filter.start_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_date.value").Optional(),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate startDate
	if data.Filters.StartDate.Value != "" && data.Filters.StartDate.Value != "all" {
		_, err := time.Parse("2006-01-02", data.Filters.StartDate.Value)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"start_date": []string{"start_date must be a valid date in format YYYY-MM-DD"},
			}
		}
	}

	return nil
}
