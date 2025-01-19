package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: EmailsUpdateStarRequestParams
 */
type EmailsUpdateStarRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *EmailsUpdateStarRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: EmailsUpdateStarRequestBody
 */
type EmailsUpdateStarRequestBody struct {
	IsStarred       string `json:"isStarred" openapi:"example:true;required"`
	IsStarredAsBool bool   `json:"-" openapi:"ignored"`
}

func (data *EmailsUpdateStarRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"isStarred": govalidity.New("isStarred").In([]string{"true", "false"}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	if data.IsStarred == "true" {
		data.IsStarredAsBool = true
	} else {
		data.IsStarredAsBool = false
	}

	return nil
}
