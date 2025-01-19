package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: QuizzesCreateStartRequestBody
 */
type QuizzesCreateStartRequestBody struct {
	QuizID   int64  `json:"quizId" openapi:"example:1;required;"`
	Password string `json:"password" openapi:"example:password;required;"`
	UserID   int64  `json:"-" openapi:"ignored"`
}

func (data *QuizzesCreateStartRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"quizId":   govalidity.New("quizId").Int().Min(1).Required(),
		"password": govalidity.New("password"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
