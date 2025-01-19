package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"net/http"
)

/*
 * @apiDefine: QuizzesUpdateEndRequestParams
 */
type QuizzesUpdateEndRequestParams struct {
	QuizID int `json:"quizId,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *QuizzesUpdateEndRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"quizId": govalidity.New("quizId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: QuizzesUpdateEndRequestBody
 */
type QuizzesUpdateEndRequestBody struct {
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"-"`
}

func (data *QuizzesUpdateEndRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
