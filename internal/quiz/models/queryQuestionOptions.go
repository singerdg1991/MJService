package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: QuizzesQueryQuestionOptionsFilterType
 */
type QuizzesQueryQuestionOptionsFilterType struct {
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsSortValue
 */
type QuizzesQueryQuestionOptionsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsSortType
 */
type QuizzesQueryQuestionOptionsSortType struct {
	CreatedAt QuizzesQueryQuestionOptionsSortValue `json:"created_at,omitempty" openapi:"$ref:QuizzesQueryQuestionOptionsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: QuizzesQueryQuestionOptionsRequestParams
 */
type QuizzesQueryQuestionOptionsRequestParams struct {
	ID      int                                   `json:"id,string,omitempty" openapi:"example:1"`
	QuizID  int64                                 `json:"quizId,string,omitempty" openapi:"example:1"`
	Page    int                                   `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                                   `json:"limit,string,omitempty" openapi:"example:10"`
	Filters QuizzesQueryQuestionOptionsFilterType `json:"filters,omitempty" openapi:"$ref:QuizzesQueryQuestionOptionsFilterType;in:query"`
	Sorts   QuizzesQueryQuestionOptionsSortType   `json:"sorts,omitempty" openapi:"$ref:QuizzesQueryQuestionOptionsSortType;in:query"`
}

func (data *QuizzesQueryQuestionOptionsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"quizId": govalidity.New("quizId").Int().Optional(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
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

	return nil
}
