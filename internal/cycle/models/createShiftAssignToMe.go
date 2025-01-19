package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateShiftAssignToMeRequestBody
 */
type CyclesCreateShiftAssignToMeRequestBody struct {
	CycleID    int                           `json:"cycleId" openapi:"example:1;required;"`
	StaffID    int                           `json:"staffId" openapi:"example:1;required;"`
	Date       string                        `json:"date" openapi:"example:2021-01-01;required;"`
	ShiftName  string                        `json:"shiftName" openapi:"example:Morning;required;"`
	Comment    *string                       `json:"comment" openapi:"example:This is a comment;"`
	DateAsDate *time.Time                    `json:"-" openapi:"ignored"`
	Staff      *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateShiftAssignToMeRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId":   govalidity.New("cycleId").Int().Min(1).Required(),
		"staffId":   govalidity.New("staffId").Int().Min(1).Required(),
		"date":      govalidity.New("date").Required(),
		"shiftName": govalidity.New("shiftName").In([]string{shifts.MorningShift, shifts.EveningShift, shifts.NightShift}).Required(),
		"comment":   govalidity.New("comment").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"date": []string{"date should have format YYYY-MM-DD"},
		}
	}

	// Reset time to 00:00:00 UTC because we only need the date
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	data.DateAsDate = &date

	return nil
}
