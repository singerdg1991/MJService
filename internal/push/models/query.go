package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: PushFilterType
 */
type PushFilterType struct {
}

/*
 * @apiDefine: PushSortValue
 */
type PushSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: PushSortType
 */
type PushSortType struct {
	CreatedAt PushSortValue `json:"created_at,omitempty" openapi:"$ref:PushSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: PushesQueryRequestParams
 */
type PushesQueryRequestParams struct {
	ID      int            `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int            `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int            `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters PushFilterType `json:"filters,omitempty" openapi:"$ref:PushFilterType;in:query"`
	Sorts   PushSortType   `json:"sorts,omitempty" openapi:"$ref:PushSortType;in:query"`
}

func (data *PushesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"userId":  govalidity.New("userId").Int().Optional(),
		"page":    govalidity.New("page").Int().Default("1"),
		"limit":   govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{},
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

	return nil
}
