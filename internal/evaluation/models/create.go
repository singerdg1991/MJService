package models

import (
	"github.com/hoitek/Maja-Service/internal/evaluation/constants"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: EvaluationsCreateRequestBody
 */
type EvaluationsCreateRequestBody struct {
	StaffID        int    `json:"staffId" openapi:"example:1;nullable;pattern:^[0-9]+$;in:body"`
	EvaluationType string `json:"evaluationType" openapi:"example:grace;required;"`
	Title          string `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Description    string `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *EvaluationsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId": govalidity.New("staffId").Int().Required(),
		"evaluationType": govalidity.New("evaluationType").In([]string{
			constants.EVALUATION_TYPE_GRACE,
			constants.EVALUATION_TYPE_WARNING,
			constants.EVALUATION_TYPE_ATTENTION,
		}),
		"title":       govalidity.New("title").MinMaxLength(3, 25).Required(),
		"description": govalidity.New("description").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
