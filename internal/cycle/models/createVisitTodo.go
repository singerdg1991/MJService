package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
)

/*
 * @apiDefine: CyclesCreateVisitTodoRequestBody
 */
type CyclesCreateVisitTodoRequestBody struct {
	CyclePickupShiftID int                             `json:"cyclePickupShiftId" openapi:"example:1;required;"`
	Title              string                          `json:"title" openapi:"example:title;required"`
	TimeValue          string                          `json:"timeValue" openapi:"example:00:00:00"`
	DateValue          string                          `json:"dateValue" openapi:"example:2021-01-01"`
	Description        *string                         `json:"description" openapi:"example:description"`
	TimeValueAsTime    *time.Time                      `json:"-" openapi:"ignored"`
	DateValueAsDate    *time.Time                      `json:"-" openapi:"ignored"`
	AuthenticatedUser  *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of an HTTP request against a predefined schema.
//
// It takes an HTTP request as a parameter and checks its body against a schema defined for the CyclesCreateVisitTodoRequestBody struct.
// r *http.Request: The HTTP request to be validated.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors that occurred.
func (data *CyclesCreateVisitTodoRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cyclePickupShiftId": govalidity.New("cyclePickupShiftId").Int().Min(1).Required(),
		"title":              govalidity.New("title").Required(),
		"timeValue":          govalidity.New("timeValue").Required(),
		"dateValue":          govalidity.New("dateValue").Required(),
		"description":        govalidity.New("description").Optional(),
	}

	// Check if the request body has error
	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate the timeValue
	t, err := time.Parse("15:04:05", data.TimeValue)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"timeValue": []string{
				"Invalid time format. Please use HH:MM:SS format",
			},
		}
	}
	data.TimeValueAsTime = &t

	// Validate the dateValue
	d, err := time.Parse("2006-01-02", data.DateValue)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"dateValue": []string{
				"Invalid date format. Please use YYYY-MM-DD format",
			},
		}
	}
	data.DateValueAsDate = &d

	return nil
}
