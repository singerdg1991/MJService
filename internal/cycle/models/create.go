package models

import (
	"fmt"
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: CyclesCreateRequestBody
 */
type CyclesCreateRequestBody struct {
	SectionID              int        `json:"sectionId" openapi:"example:1;required;"`
	StartDate              string     `json:"start_date" openapi:"example:2021-01-01;required;"`
	EndDate                *string    `json:"end_date" openapi:"example:2021-01-01;required;"`
	PeriodLength           *string    `json:"periodLength" openapi:"example:oneWeek;required;"`
	ShiftMorningStartTime  string     `json:"shiftMorningStartTime" openapi:"example:08:00;required;"`
	ShiftMorningEndTime    string     `json:"shiftMorningEndTime" openapi:"example:16:00;required;"`
	ShiftEveningStartTime  string     `json:"shiftEveningStartTime" openapi:"example:16:00;required;"`
	ShiftEveningEndTime    string     `json:"shiftEveningEndTime" openapi:"example:00:00;required;"`
	ShiftNightStartTime    string     `json:"shiftNightStartTime" openapi:"example:00:00;required;"`
	ShiftNightEndTime      string     `json:"shiftNightEndTime" openapi:"example:08:00;required;"`
	FreezePeriodDate       string     `json:"freeze_period_date" openapi:"example:2021-01-01;required;"`
	WishDays               int        `json:"wishDays" openapi:"example:1;required;"`
	Name                   string     `json:"name" openapi:"ignored"`
	StartDateAsDate        *time.Time `json:"-" openapi:"ignored"`
	EndDateAsDate          *time.Time `json:"-" openapi:"ignored"`
	FreezePeriodDateAsDate *time.Time `json:"-" openapi:"ignored"`
}

// ValidateBody validates the request body of the CyclesCreateRequestBody.
//
// It takes a http.Request as a parameter and returns a govalidity.ValidityResponseErrors.
func (data *CyclesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"sectionId":             govalidity.New("sectionId").Int().Min(1).Required(),
		"start_date":            govalidity.New("start_date").Required(),
		"end_date":              govalidity.New("end_date"),
		"periodLength":          govalidity.New("periodLength"),
		"shiftMorningStartTime": govalidity.New("shiftMorningStartTime").Required(),
		"shiftMorningEndTime":   govalidity.New("shiftMorningEndTime").Required(),
		"shiftEveningStartTime": govalidity.New("shiftEveningStartTime").Required(),
		"shiftEveningEndTime":   govalidity.New("shiftEveningEndTime").Required(),
		"shiftNightStartTime":   govalidity.New("shiftNightStartTime").Required(),
		"shiftNightEndTime":     govalidity.New("shiftNightEndTime").Required(),
		"freeze_period_date":    govalidity.New("freeze_period_date").Required(),
		"wishDays":              govalidity.New("wishDays").Int().Min(0).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate startDate
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"Start date is invalid"},
		}
	}
	data.StartDateAsDate = &startDate

	// Check start date is in the future
	if startDate.Before(time.Now()) {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"Start date is in the past"},
		}
	}

	// Validate endDate
	if data.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *data.EndDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"End date is invalid"},
			}
		}
		data.EndDateAsDate = &endDate

		// Check end date is in the future
		if endDate.Before(startDate) {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"End date is in the past"},
			}
		}
	}

	// Validate freezePeriodDate
	freezePeriodDate, err := time.Parse("2006-01-02", data.FreezePeriodDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"freeze_period_date": []string{"Freeze period date is invalid"},
		}
	}
	data.FreezePeriodDateAsDate = &freezePeriodDate

	// Check periodLength and endDate
	if data.PeriodLength == nil && data.EndDate == nil {
		return govalidity.ValidityResponseErrors{
			"periodLength": []string{"Period length or end date is required"},
		}
	}

	// Validate periodLength
	if data.PeriodLength != nil {
		if data.EndDate != nil {
			return govalidity.ValidityResponseErrors{
				"periodLength": []string{"Period length and end date can not be used together"},
			}
		}
		if *data.PeriodLength != constants.CYCLE_PERIOD_LENGTH_ONE_WEEK && *data.PeriodLength != constants.CYCLE_PERIOD_LENGTH_TWO_WEEKS && *data.PeriodLength != constants.CYCLE_PERIOD_LENGTH_THREE_WEEKS {
			return govalidity.ValidityResponseErrors{
				"periodLength": []string{fmt.Sprintf("Period length must be one of %s, %s, %s", constants.CYCLE_PERIOD_LENGTH_ONE_WEEK, constants.CYCLE_PERIOD_LENGTH_TWO_WEEKS, constants.CYCLE_PERIOD_LENGTH_THREE_WEEKS)},
			}
		}

		// Calculate end date based on period length
		oneWeek := 7
		switch *data.PeriodLength {
		case constants.CYCLE_PERIOD_LENGTH_ONE_WEEK:
			// add one week to start date
			endDate := data.StartDateAsDate.AddDate(0, 0, oneWeek)
			data.EndDateAsDate = &endDate
		case constants.CYCLE_PERIOD_LENGTH_TWO_WEEKS:
			// add two weeks to start date
			endDate := data.StartDateAsDate.AddDate(0, 0, 2*oneWeek)
			data.EndDateAsDate = &endDate
		case constants.CYCLE_PERIOD_LENGTH_THREE_WEEKS:
			// add three weeks to start date
			endDate := data.StartDateAsDate.AddDate(0, 0, 3*oneWeek)
			data.EndDateAsDate = &endDate
		}
	}

	// Validate endDate
	if data.EndDate != nil {
		// Check if period length is set
		if data.PeriodLength != nil {
			return govalidity.ValidityResponseErrors{
				"end_date": []string{"Period length and end date can not be used together"},
			}
		}
	}

	// Check if end date is before start date
	if data.EndDateAsDate.Before(*data.StartDateAsDate) {
		return govalidity.ValidityResponseErrors{
			"end_date": []string{"End date must be after start date"},
		}
	}

	// Check if end date is in the past
	if data.EndDateAsDate.Before(time.Now()) {
		return govalidity.ValidityResponseErrors{
			"end_date": []string{"End date must be in the future"},
		}
	}

	// Validate freezePeriodDate
	if data.FreezePeriodDateAsDate.Before(*data.StartDateAsDate) {
		return govalidity.ValidityResponseErrors{
			"freeze_period_date": []string{"Freeze period date must be after start date"},
		}
	}
	if data.FreezePeriodDateAsDate.After(*data.EndDateAsDate) || data.FreezePeriodDateAsDate.Equal(*data.EndDateAsDate) {
		return govalidity.ValidityResponseErrors{
			"freeze_period_date": []string{"Freeze period date must be before end date"},
		}
	}

	// Validate shiftMorningStartTime
	_, err = time.Parse("15:04", data.ShiftMorningStartTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftMorningStartTime": []string{"Shift morning start time is invalid"},
		}
	}

	// Validate shiftMorningEndTime
	_, err = time.Parse("15:04", data.ShiftMorningEndTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftMorningEndTime": []string{"Shift morning end time is invalid"},
		}
	}

	// Validate shiftEveningStartTime
	_, err = time.Parse("15:04", data.ShiftEveningStartTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftEveningStartTime": []string{"Shift evening start time is invalid"},
		}
	}

	// Validate shiftEveningEndTime
	_, err = time.Parse("15:04", data.ShiftEveningEndTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftEveningEndTime": []string{"Shift evening end time is invalid"},
		}
	}

	// Validate shiftNightStartTime
	_, err = time.Parse("15:04", data.ShiftNightStartTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftNightStartTime": []string{"Shift night start time is invalid"},
		}
	}

	// Validate shiftNightEndTime
	_, err = time.Parse("15:04", data.ShiftNightEndTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"shiftNightEndTime": []string{"Shift night end time is invalid"},
		}
	}

	return nil
}
