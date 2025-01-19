package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Maja-Service/internal/city/domain"
)

/*
 * @apiDefine: CityFilterType
 */
type CityFilterType struct {
	Name      filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: CitiesQueryRequestParams
 */
type CitiesQueryRequestParams struct {
	ID      int            `json:"id,string,omitempty" openapi:"example:1"`
	Page    int            `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters CityFilterType `json:"filters,omitempty" openapi:"$ref:CityFilterType;in:query"`
}

func (data *CitiesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
 * @apiDefine: CitiesQueryResponseDataItem
 */
type CitiesQueryResponseDataItem struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: CitiesQueryResponseData
 */
type CitiesQueryResponseData struct {
	Limit      int                           `json:"limit" openapi:"example:10"`
	Offset     int                           `json:"offset" openapi:"example:0"`
	Page       int                           `json:"page" openapi:"example:1"`
	TotalRows  int                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                           `json:"totalPages" openapi:"example:1"`
	Items      []CitiesQueryResponseDataItem `json:"items" openapi:"$ref:CitiesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: CitiesQueryResponse
 */
type CitiesQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       CitiesQueryResponseData `json:"data" openapi:"$ref:CitiesQueryResponseData"`
}

/*
 * @apiDefine: CitiesQueryNotFoundResponse
 */
type CitiesQueryNotFoundResponse struct {
	Cities []domain.City `json:"cities" openapi:"$ref:City;type:array"`
}
