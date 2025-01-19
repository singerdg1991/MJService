package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: CyclesShiftCustomerHomeKeyNotReleaseRequestBody
 */
type CyclesShiftCustomerHomeKeyNotReleaseRequestBody struct {
	ShiftID           int                             `json:"shiftId" openapi:"example:1;required;"`
	KeyNo             string                          `json:"keyNo" openapi:"example:1;required;"`
	Reason            *string                         `json:"reason" openapi:"example:reason;"`
	Status            string                          `json:"-" openapi:"ignored"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the request body of CyclesShiftCustomerHomeKeyNotReleaseRequestBody.
//
// It takes an http.Request as a parameter and checks if the request body matches the defined schema.
// If the validation fails, it returns a govalidity.ValidityResponseErrors with the errors.
// Otherwise, it sets the status of the request body to constants.SHIFT_CUSTOMER_HOME_KEY_STATUS_NOT_RELEASED
// and returns nil.
//
// Parameters:
// - r: The http.Request to validate.
//
// Returns:
// - govalidity.ValidityResponseErrors: The errors encountered during validation, if any.
func (data *CyclesShiftCustomerHomeKeyNotReleaseRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"shiftId": govalidity.New("shiftId").Int().Min(1).Required(),
		"keyNo":   govalidity.New("KeyNo").Required(),
		"reason":  govalidity.New("Status").Optional(),
	}

	// Validate the request body
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Set status to "released"
	data.Status = constants.SHIFT_CUSTOMER_HOME_KEY_STATUS_NOT_RELEASED

	return nil
}
