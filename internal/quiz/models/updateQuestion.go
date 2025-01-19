package models

import (
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: QuizzesUpdateQuestionRequestParams
 */
type QuizzesUpdateQuestionRequestParams struct {
	QuestionID int `json:"questionId,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *QuizzesUpdateQuestionRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"questionId": govalidity.New("questionId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
