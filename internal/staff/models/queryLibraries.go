package models

import (
	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsQueryLibrariesSortValue
 */
type StaffsQueryLibrariesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: StaffsQueryLibrariesSortType
 */
type StaffsQueryLibrariesSortType struct {
	Title     StaffsQueryLibrariesSortValue `json:"title,omitempty" openapi:"$ref:StaffsQueryLibrariesSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt StaffsQueryLibrariesSortValue `json:"created_at,omitempty" openapi:"$ref:StaffsQueryLibrariesSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: StaffsQueryLibrariesFilterType
 */
type StaffsQueryLibrariesFilterType struct {
	Title filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"title\"}}"`
}

/*
 * @apiDefine: StaffsQueryLibrariesRequestParams
 */
type StaffsQueryLibrariesRequestParams struct {
	ID      int                            `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int                            `json:"userId,string,omitempty" openapi:"example:1;required"`
	Page    int                            `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters StaffsQueryLibrariesFilterType `json:"filters,omitempty" openapi:"$ref:StaffsQueryLibrariesFilterType;in:query"`
	Sorts   StaffsQueryLibrariesSortType   `json:"sorts,omitempty" openapi:"$ref:StaffsQueryLibrariesSortType;in:query"`
}

func (data *StaffsQueryLibrariesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"userId": govalidity.New("userId").Int().Required(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
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
