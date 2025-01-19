package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: QuizzesCreateQuestionRequestBodyOption
 */
type QuizzesCreateQuestionRequestBodyOption struct {
	Title string `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Score int    `json:"score" openapi:"example:0;required;"`
}

/*
 * @apiDefine: QuizzesCreateQuestionRequestBody
 */
type QuizzesCreateQuestionRequestBody struct {
	QuizID      int64                                     `json:"quizId" openapi:"example:1;required;"`
	Title       string                                    `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Description *string                                   `json:"description" openapi:"example:description;required;"`
	Options     []*QuizzesCreateQuestionRequestBodyOption `json:"options" openapi:"$ref:QuizzesCreateQuestionRequestBodyOption;type:array;required;"`
}

func (data *QuizzesCreateQuestionRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"quizId":      govalidity.New("quizId").Int().Min(1).Required(),
		"title":       govalidity.New("title").MinMaxLength(3, 25).Required(),
		"description": govalidity.New("description"),
		"options":     govalidity.New("options"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate options
	if len(data.Options) > 0 {
		for _, option := range data.Options {
			if option == nil {
				continue
			}
			if option.Title == "" {
				return govalidity.ValidityResponseErrors{
					"options.title": []string{"title is required"},
				}
			}
			if option.Score < 0 {
				return govalidity.ValidityResponseErrors{
					"options.score": []string{"score must be greater than or equal to 0"},
				}
			}
		}
	}

	return nil
}
