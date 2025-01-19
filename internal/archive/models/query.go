package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: ArchiveFilterType
 */
type ArchiveFilterType struct {
	Title       filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Description filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: ArchiveSortValue
 */
type ArchiveSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: ArchiveSortType
 */
type ArchiveSortType struct {
	Title       ArchiveSortValue `json:"title,omitempty" openapi:"$ref:ArchiveSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	Description ArchiveSortValue `json:"description,omitempty" openapi:"$ref:ArchiveSortValue;example:{\"description\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: ArchivesQueryRequestParams
 */
type ArchivesQueryRequestParams struct {
	ID      int               `json:"id,string,omitempty" openapi:"example:1"`
	Page    int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters ArchiveFilterType `json:"filters,omitempty" openapi:"$ref:ArchiveFilterType;in:query"`
	Sorts   ArchiveSortType   `json:"sorts,omitempty" openapi:"$ref:ArchiveSortType;in:query"`
}

func (data *ArchivesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"title": govalidity.Schema{
				"op": govalidity.New("sort.title.op"),
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
