package models

import (
	"fmt"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/section/constants"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: SectionSortValue
 */
type SectionSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: SectionSortType
 */
type SectionSortType struct {
	CreatedAt SectionSortValue `json:"created_at,omitempty" openapi:"$ref:SectionSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: SectionFilterType
 */
type SectionFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: SectionsQueryRequestParams
 */
type SectionsQueryRequestParams struct {
	ID       int               `json:"id,string,omitempty" openapi:"example:1"`
	ParentID int               `json:"parentId,string,omitempty" openapi:"example:1"`
	Page     int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit    int               `json:"limit,string,omitempty" openapi:"example:10"`
	Type     string            `json:"type,omitempty" openapi:"example:all"`
	Filters  SectionFilterType `json:"filters,omitempty" openapi:"$ref:SectionFilterType;in:query"`
	Sorts    SectionSortType   `json:"sorts,omitempty" openapi:"$ref:SectionSortType;in:query"`
}

func (data *SectionsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Int().Optional(),
		"parentId": govalidity.New("parentId").Int().Optional(),
		"page":     govalidity.New("page").Int().Default("1"),
		"limit":    govalidity.New("limit").Int().Default("10"),
		"type":     govalidity.New("type"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"parentId": govalidity.Schema{
				"op":    govalidity.New("filter.parentId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.parentId.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
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

	// Check type
	if data.Type == "" {
		data.Type = constants.SECTION_TYPE_ALL
	}
	if data.Type != constants.SECTION_TYPE_ALL && data.Type != constants.SECTION_TYPE_PARENT && data.Type != constants.SECTION_TYPE_CHILDREN {
		return govalidity.ValidityResponseErrors{
			"type": []string{fmt.Sprintf("The type must be one of %s, %s, %s", constants.SECTION_TYPE_ALL, constants.SECTION_TYPE_PARENT, constants.SECTION_TYPE_CHILDREN)},
		}
	}

	return nil
}
