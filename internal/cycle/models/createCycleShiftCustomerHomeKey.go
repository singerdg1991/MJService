package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: CyclesShiftCustomerHomeKeyRequestBody
 */
type CyclesShiftCustomerHomeKeyRequestBody struct {
	ShiftID           int                             `json:"shiftId" openapi:"example:1;required;"`
	KeyNo             string                          `json:"keyNo" openapi:"example:1;required;"`
	Status            string                          `json:"status" openapi:"example:accepted;required;"`
	Reason            *string                         `json:"reason" openapi:"example:reason;"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the request body of the CyclesShiftCustomerHomeKeyRequestBody.
//
// The function takes an http.Request as a parameter and checks if the request body matches the defined schema.
// It returns a govalidity.ValidityResponseErrors if the validation fails, otherwise it returns nil.
func (data *CyclesShiftCustomerHomeKeyRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"shiftId": govalidity.New("shiftId").Int().Min(1).Required(),
		"keyNo":   govalidity.New("KeyNo").Required(),
		"status": govalidity.New("Status").In([]string{
			constants.SHIFT_CUSTOMER_HOME_KEY_STATUS_ACCEPTED,
			constants.SHIFT_CUSTOMER_HOME_KEY_STATUS_REJECTED,
		}),
		"reason": govalidity.New("reason").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
