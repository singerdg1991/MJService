package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: QuizParticipantFilterType
 */
type QuizParticipantFilterType struct {
	Title     filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: QuizParticipantSortValue
 */
type QuizParticipantSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: QuizParticipantSortType
 */
type QuizParticipantSortType struct {
	Title     QuizParticipantSortValue `json:"title,omitempty" openapi:"$ref:QuizParticipantSortValue;example:{\"title\":{\"op\":\"asc\"}}"`
	CreatedAt QuizParticipantSortValue `json:"created_at,omitempty" openapi:"$ref:QuizParticipantSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: QuizzesQueryParticipantsRequestParams
 */
type QuizzesQueryParticipantsRequestParams struct {
	ID      int                       `json:"id,string,omitempty" openapi:"example:1"`
	QuizID  int64                     `json:"quizId,string,omitempty" openapi:"example:1"`
	UserID  int64                     `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int                       `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                       `json:"limit,string,omitempty" openapi:"example:10"`
	Filters QuizParticipantFilterType `json:"filters,omitempty" openapi:"$ref:QuizParticipantFilterType;in:query"`
	Sorts   QuizParticipantSortType   `json:"sorts,omitempty" openapi:"$ref:QuizParticipantSortType;in:query"`
}

func (data *QuizzesQueryParticipantsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":     govalidity.New("id").Int().Optional(),
		"quizId": govalidity.New("quizId").Int().Optional(),
		"userId": govalidity.New("userId").Int().Optional(),
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
