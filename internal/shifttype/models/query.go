package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Maja-Service/internal/shifttype/domain"
)

/*
 * @apiDefine: ShiftTypeFilterType
 */
type ShiftTypeFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: ShiftTypesQueryRequestParams
 */
type ShiftTypesQueryRequestParams struct {
	ID      int                 `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                 `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                 `json:"limit,string,omitempty" openapi:"example:10"`
	Filters ShiftTypeFilterType `json:"filters,omitempty" openapi:"$ref:ShiftTypeFilterType;in:query"`
}

func (data *ShiftTypesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
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

/*
 * @apiDefine: ShiftTypesQueryResponseDataItem
 */
type ShiftTypesQueryResponseDataItem struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: ShiftTypesQueryResponseData
 */
type ShiftTypesQueryResponseData struct {
	Limit      int                               `json:"limit" openapi:"example:10"`
	Offset     int                               `json:"offset" openapi:"example:0"`
	Page       int                               `json:"page" openapi:"example:1"`
	TotalRows  int                               `json:"totalRows" openapi:"example:1"`
	TotalPages int                               `json:"totalPages" openapi:"example:1"`
	Items      []ShiftTypesQueryResponseDataItem `json:"items" openapi:"$ref:ShiftTypesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ShiftTypesQueryResponse
 */
type ShiftTypesQueryResponse struct {
	StatusCode int                         `json:"statusCode" openapi:"example:200"`
	Data       ShiftTypesQueryResponseData `json:"data" openapi:"$ref:ShiftTypesQueryResponseData"`
}

/*
 * @apiDefine: ShiftTypesQueryNotFoundResponse
 */
type ShiftTypesQueryNotFoundResponse struct {
	ShiftTypes []domain.ShiftType `json:"ShiftTypes" openapi:"$ref:ShiftType;type:array"`
}
