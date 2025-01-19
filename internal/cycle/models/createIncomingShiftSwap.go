package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateIncomingShiftSwapRequestBody
 */
type CyclesCreateIncomingShiftSwapRequestBody struct {
	CycleID          int                           `json:"cycleId" openapi:"example:1;required;"`
	SourceStaffID    int                           `json:"sourceStaffId" openapi:"example:1;required;"`
	SourceDate       string                        `json:"sourceDate" openapi:"example:2021-01-01;required;"`
	SourceShiftName  string                        `json:"sourceShiftName" openapi:"example:Morning;required;"`
	TargetStaffID    int                           `json:"targetStaffId" openapi:"example:1;required;"`
	TargetDate       string                        `json:"targetDate" openapi:"example:2021-01-01;required;"`
	TargetShiftName  string                        `json:"targetShiftName" openapi:"example:Morning;required;"`
	Comment          *string                       `json:"comment" openapi:"example:This is a comment;"`
	SourceDateAsDate *time.Time                    `json:"-" openapi:"ignored"`
	TargetDateAsDate *time.Time                    `json:"-" openapi:"ignored"`
	SourceStaff      *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
	TargetStaff      *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of an HTTP request against a predefined schema.
//
// The function takes an HTTP request as a parameter and checks its body against a schema defined for the CyclesCreateIncomingShiftSwapRequestBody struct.
// It returns a govalidity.ValidityResponseErrors object containing any validation errors that occurred.
func (data *CyclesCreateIncomingShiftSwapRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId":         govalidity.New("cycleId").Int().Min(1).Required(),
		"sourceStaffId":   govalidity.New("sourceStaffId").Int().Min(1).Required(),
		"sourceDate":      govalidity.New("sourceDate").Required(),
		"sourceShiftName": govalidity.New("sourceShiftName").In([]string{shifts.MorningShift, shifts.EveningShift, shifts.NightShift}).Required(),
		"targetStaffId":   govalidity.New("targetStaffId").Int().Min(1).Required(),
		"targetDate":      govalidity.New("targetDate").Required(),
		"targetShiftName": govalidity.New("targetShiftName").In([]string{shifts.MorningShift, shifts.EveningShift, shifts.NightShift}).Required(),
		"comment":         govalidity.New("comment").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	sourceDate, err := time.Parse("2006-01-02", data.SourceDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"sourceDate": []string{"sourceDate should have format YYYY-MM-DD"},
		}
	}

	// Reset time to 00:00:00 UTC because we only need the date
	sourceDate = time.Date(sourceDate.Year(), sourceDate.Month(), sourceDate.Day(), 0, 0, 0, 0, time.UTC)
	data.SourceDateAsDate = &sourceDate

	// Validate datetime
	targetDate, err := time.Parse("2006-01-02", data.TargetDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"targetDate": []string{"targetDate should have format YYYY-MM-DD"},
		}
	}

	// Reset time to 00:00:00 UTC because we only need the date
	targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	data.TargetDateAsDate = &targetDate

	return nil
}
