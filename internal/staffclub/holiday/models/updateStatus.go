package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/constants"
	"net/http"
)

/*
 * @apiDefine: HolidaysUpdateStatusRequestParams
 */
type HolidaysUpdateStatusRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *HolidaysUpdateStatusRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: HolidaysUpdateStatusRequestBody
 */
type HolidaysUpdateStatusRequestBody struct {
	Status            string                         `json:"status" openapi:"example:pending;required"`
	RejectedReason    *string                        `json:"rejectedReason" openapi:"example:reason;nullable"`
	AuthenticatedUser sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *HolidaysUpdateStatusRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"status":         govalidity.New("status").In([]string{constants.HOLIDAY_STATUS_PENDING, constants.HOLIDAY_STATUS_ACCEPTED, constants.HOLIDAY_STATUS_REJECTED}).Required(),
		"rejectedReason": govalidity.New("rejectedReason"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
