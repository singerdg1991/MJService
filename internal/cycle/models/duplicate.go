package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesDuplicateRequestBody
 */
type CyclesDuplicateRequestBody struct {
	Cycle           *domain.Cycle            `json:"-" openapi:"ignored"`
	CycleStaffTypes []*domain.CycleStaffType `json:"-" openapi:"ignored"`
	CycleID         int64                    `json:"cycleId" openapi:"example:1;required;"`
}

// ValidateBody validates the body of the CyclesDuplicateRequestBody against the defined schema.
//
// It takes an http.Request object as a parameter and checks if the request body matches the defined schema.
// The schema includes a "cycleId" field which is required and must be an integer with a minimum value of 1.
// If the validation fails, it returns a govalidity.ValidityResponseErrors object with the corresponding error messages.
// Otherwise, it returns nil.
func (data *CyclesDuplicateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId": govalidity.New("cycleId").Int().Min(1).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
