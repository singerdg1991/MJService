package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesUpdateNextStaffTypeRequestParams
 */
type CyclesUpdateNextStaffTypeRequestParams struct {
	CurrentCycleID int `json:"currentcycleid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

// ValidateParams validates the params of an HTTP request against a predefined schema.
//
// It takes a govalidity.Params object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesUpdateNextStaffTypeRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"currentcycleid": govalidity.New("currentcycleid").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypeRequestBody
 */
type CyclesUpdateNextStaffTypeRequestBody struct {
	ShiftName        string                         `json:"shiftName" openapi:"example:morning;required;"`
	DateTime         string                         `json:"datetime" openapi:"example:2021-01-01;required;"`
	StartHour        string                         `json:"startHour" openapi:"example:00:00"`
	EndHour          string                         `json:"endHour" openapi:"example:00:00"`
	NeededStaffCount int                            `json:"neededStaffCount" openapi:"example:1;required;"`
	RoleID           int                            `json:"roleId" openapi:"example:1;required;"`
	DateTimeAsDate   *time.Time                     `json:"-" openapi:"ignored"`
	Role             *domain.CycleNextStaffTypeRole `json:"-" openapi:"ignored"`
	StartHourAsTime  *time.Time                     `json:"-" openapi:"ignored"`
	EndHourAsTime    *time.Time                     `json:"-" openapi:"ignored"`
}

func (data *CyclesUpdateNextStaffTypeRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"shiftName":        govalidity.New("shiftName").In([]string{shifts.MorningShift, shifts.EveningShift, shifts.NightShift}).Required(),
		"datetime":         govalidity.New("datetime").Required(),
		"neededStaffCount": govalidity.New("neededStaffCount").Int().Min(0).Required(),
		"roleId":           govalidity.New("roleId").Int().Min(1).Required(),
		"startHour":        govalidity.New("startHour").Required(),
		"endHour":          govalidity.New("endHour").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	datetime, err := time.Parse("2006-01-02", data.DateTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"datetime": []string{"datetime should have format YYYY-MM-DD"},
		}
	}

	// Reset time to 00:00:00 UTC because we only need the date
	datetime = time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, time.UTC)
	data.DateTimeAsDate = &datetime

	// Validate startHour
	startTime, err := time.Parse("15:04", data.StartHour)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"startHour": []string{"startHour should have format HH:MM"},
		}
	}
	data.StartHourAsTime = &startTime

	// Validate endHour
	endTime, err := time.Parse("15:04", data.EndHour)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"endHour": []string{"endHour should have format HH:MM"},
		}
	}
	data.EndHourAsTime = &endTime

	return nil
}
