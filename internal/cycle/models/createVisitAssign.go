package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateVisitAssignRequestBody
 */
type CyclesCreateVisitAssignRequestBody struct {
	CyclePickupShiftID int                           `json:"cyclePickupShiftId" openapi:"example:1;required;"`
	StaffID            int                           `json:"staffId" openapi:"example:1;required;"`
	Comment            *string                       `json:"comment" openapi:"example:This is a comment;"`
	Staff              *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateVisitAssignRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cyclePickupShiftId": govalidity.New("cyclePickupShiftId").Int().Min(1).Required(),
		"staffId":            govalidity.New("staffId").Int().Min(1).Required(),
		"comment":            govalidity.New("comment").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
