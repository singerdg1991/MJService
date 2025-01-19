package models

import (
	"net/http"
	"strconv"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
)

/*
 * @apiDefine: CyclesCreateVisitStartRequestBody
 */
type CyclesCreateVisitStartRequestBody struct {
	CyclePickupShiftID int                             `json:"cyclePickupShiftId" openapi:"example:1;required;"`
	StartKilometer     string                          `json:"startKilometer" openapi:"example:100;"`
	AuthenticatedUser  *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request against a predefined schema.
//
// It takes an http.Request object as a parameter and returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateVisitStartRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cyclePickupShiftId": govalidity.New("cyclePickupShiftId").Int().Min(1).Required(),
		"startKilometer":     govalidity.New("startKilometer").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Check if startKilometer is number
	if _, err := strconv.ParseFloat(data.StartKilometer, 64); err != nil {
		return govalidity.ValidityResponseErrors{
			"startKilometer": []string{"startKilometer must be a number"},
		}
	}

	return nil
}
