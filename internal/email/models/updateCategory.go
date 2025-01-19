package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/email/constants"
	"net/http"
)

/*
 * @apiDefine: EmailsUpdateCategoryRequestParams
 */
type EmailsUpdateCategoryRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *EmailsUpdateCategoryRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: EmailsUpdateCategoryRequestBody
 */
type EmailsUpdateCategoryRequestBody struct {
	Category string `json:"category" openapi:"example:outbox;required"`
}

func (data *EmailsUpdateCategoryRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"category": govalidity.New("category").In([]string{
			constants.EMAIL_CATEGORY_OUTBOX,
			constants.EMAIL_CATEGORY_DRAFT,
			constants.EMAIL_CATEGORY_ARCHIVE,
			constants.EMAIL_CATEGORY_TRASH,
		}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
