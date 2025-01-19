package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: QuizFilterType
 */
type QuizFilterType struct {
	Title     filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: QuizSortValue
 */
type QuizSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: QuizSortType
 */
type QuizSortType struct {
	Title     QuizSortValue `json:"title,omitempty" openapi:"$ref:QuizSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt QuizSortValue `json:"created_at,omitempty" openapi:"$ref:QuizSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: QuizzesQueryRequestParams
 */
type QuizzesQueryRequestParams struct {
	ID                int                             `json:"id,string,omitempty" openapi:"example:1"`
	Page              int                             `json:"page,string,omitempty" openapi:"example:1"`
	Limit             int                             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters           QuizFilterType                  `json:"filters,omitempty" openapi:"$ref:QuizFilterType;in:query"`
	Sorts             QuizSortType                    `json:"sorts,omitempty" openapi:"$ref:QuizSortType;in:query"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *QuizzesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
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
