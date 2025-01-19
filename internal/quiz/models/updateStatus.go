package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/quiz/constants"
	"net/http"
)

/*
 * @apiDefine: QuizzesUpdateStatusRequestParams
 */
type QuizzesUpdateStatusRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *QuizzesUpdateStatusRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id": govalidity.New("id").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: QuizzesUpdateStatusRequestBody
 */
type QuizzesUpdateStatusRequestBody struct {
	Status string `json:"status" openapi:"example:disable;required;"`
}

func (data *QuizzesUpdateStatusRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"status": govalidity.New("status").In([]string{constants.QUIZ_STATUS_ENABLE, constants.QUIZ_STATUS_DISABLE}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
