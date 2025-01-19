package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CyclesCreateVisitSwapRequestBody
 */
type CyclesCreateVisitSwapRequestBody struct {
	SourceCyclePickupShiftID int     `json:"sourceCyclePickupShiftId" openapi:"example:1;required;"`
	TargetCyclePickupShiftID int     `json:"targetCyclePickupShiftId" openapi:"example:1;required;"`
	Comment                  *string `json:"comment" openapi:"example:This is a comment;"`
}

// ValidateBody validates the body of an HTTP request against a predefined schema.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateVisitSwapRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"sourceCyclePickupShiftId": govalidity.New("sourceCyclePickupShiftId").Int().Min(1).Required(),
		"targetCyclePickupShiftId": govalidity.New("targetCyclePickupShiftId").Int().Min(1).Required(),
		"comment":                  govalidity.New("comment").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
