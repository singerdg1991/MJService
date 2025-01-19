package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"net/http"
)

/*
 * @apiDefine: NotificationsUpdateStatusRequestParams
 */
type NotificationsUpdateStatusRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *NotificationsUpdateStatusRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: NotificationsUpdateStatusRequestBody
 */
type NotificationsUpdateStatusRequestBody struct {
	Status string `json:"status" openapi:"example:pending;required;maxLen:100;minLen:2;"`
}

func (data *NotificationsUpdateStatusRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"status": govalidity.New("status").In([]string{
			constants.NOTIFICATION_STATUS_PENDING,
			constants.NOTIFICATION_STATUS_ACCEPTED,
			constants.NOTIFICATION_STATUS_REJECTED,
		}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
