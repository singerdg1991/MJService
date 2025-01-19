package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: QuizQuestionFilterType
 */
type QuizQuestionFilterType struct {
	Title     filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: QuizQuestionSortValue
 */
type QuizQuestionSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: QuizQuestionSortType
 */
type QuizQuestionSortType struct {
	Title     QuizQuestionSortValue `json:"title,omitempty" openapi:"$ref:QuizQuestionSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt QuizQuestionSortValue `json:"created_at,omitempty" openapi:"$ref:QuizQuestionSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: QuizzesQueryQuestionsRequestParams
 */
type QuizzesQueryQuestionsRequestParams struct {
	ID      int                    `json:"id,string,omitempty" openapi:"example:1"`
	QuizID  int64                  `json:"quizId,string,omitempty" openapi:"example:1"`
	Page    int                    `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                    `json:"limit,string,omitempty" openapi:"example:10"`
	Filters QuizQuestionFilterType `json:"filters,omitempty" openapi:"$ref:QuizQuestionFilterType;in:query"`
	Sorts   QuizQuestionSortType   `json:"sorts,omitempty" openapi:"$ref:QuizQuestionSortType;in:query"`
}

func (data *QuizzesQueryQuestionsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"quizId": govalidity.New("quizId").Int().Optional(),
		"page":   govalidity.New("page").Int().Default("1"),
		"limit":  govalidity.New("limit").Int().Default("10"),
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
