package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
)

/*
 * @apiDefine: CyclesCreateVisitResumeRequestBody
 */
type CyclesCreateVisitResumeRequestBody struct {
	CyclePickupShiftID int                             `json:"cyclePickupShiftId" openapi:"example:1;required;"`
	Reason             string                          `json:"reason" openapi:"example:reason;"`
	AuthenticatedUser  *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateVisitResumeRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cyclePickupShiftId": govalidity.New("cyclePickupShiftId").Int().Min(1).Required(),
		"reason":             govalidity.New("reason").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
