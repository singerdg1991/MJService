package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: HolidayFilterType
 */
type HolidayFilterType struct {
	Title  filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Status filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"pending\"}"`
}

/*
 * @apiDefine: HolidaySortValue
 */
type HolidaySortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: HolidaySortType
 */
type HolidaySortType struct {
	ID        HolidaySortValue `json:"id,omitempty" openapi:"$ref:HolidaySortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	Title     HolidaySortValue `json:"title,omitempty" openapi:"$ref:HolidaySortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt HolidaySortValue `json:"created_at,omitempty" openapi:"$ref:HolidaySortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: HolidaysQueryRequestParams
 */
type HolidaysQueryRequestParams struct {
	ID      int               `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int               `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters HolidayFilterType `json:"filters,omitempty" openapi:"$ref:HolidayFilterType;in:query"`
	Sorts   HolidaySortType   `json:"sorts,omitempty" openapi:"$ref:HolidaySortType;in:query"`
}

func (data *HolidaysQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"userId": govalidity.New("userId").Int().Optional(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
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
