package models

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/shifts"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesUpdateNextStaffTypesRequestParams
 */
type CyclesUpdateNextStaffTypesRequestParams struct {
	CurrentCycleID int `json:"currentcycleid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

// ValidateParams validates the params of an HTTP request against a predefined schema.
//
// It takes a govalidity.Params object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesUpdateNextStaffTypesRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
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
 * @apiDefine: CyclesUpdateNextStaffTypesRequestBodyStaffType
 */
type CyclesUpdateNextStaffTypesRequestBodyStaffType struct {
	ShiftName        string                     `json:"shiftName" openapi:"example:morning;required;"`
	DateTime         string                     `json:"datetime" openapi:"example:2021-01-01;required;"`
	StartHour        string                     `json:"startHour" openapi:"example:00:00"`
	EndHour          string                     `json:"endHour" openapi:"example:00:00"`
	NeededStaffCount int                        `json:"neededStaffCount" openapi:"example:1;required;"`
	RoleID           int                        `json:"roleId" openapi:"example:1;required;"`
	DateTimeAsDate   *time.Time                 `json:"-" openapi:"ignored"`
	Role             *domain.CycleStaffTypeRole `json:"-" openapi:"ignored"`
	StartHourAsTime  *time.Time                 `json:"-" openapi:"ignored"`
	EndHourAsTime    *time.Time                 `json:"-" openapi:"ignored"`
}

/*
 * @apiDefine: CyclesUpdateNextStaffTypesRequestBody
 */
type CyclesUpdateNextStaffTypesRequestBody struct {
	StaffTypes []*CyclesUpdateNextStaffTypeRequestBody `json:"staffTypes" openapi:"$ref:CyclesUpdateNextStaffTypesRequestBodyStaffType"`
}

func (data *CyclesUpdateNextStaffTypesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffTypes": govalidity.New("staffTypes"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Check if staffTypes is empty
	if len(data.StaffTypes) == 0 {
		return govalidity.ValidityResponseErrors{
			"staffTypes": []string{"Staff types are required"},
		}
	}

	// Validate each staffType
	for _, staffType := range data.StaffTypes {
		// Validate shiftName
		if staffType.ShiftName != shifts.MorningShift && staffType.ShiftName != shifts.EveningShift && staffType.ShiftName != shifts.NightShift {
			return govalidity.ValidityResponseErrors{
				"shiftName": []string{fmt.Sprintf("shiftName can only be [%s %s %s]", shifts.MorningShift, shifts.EveningShift, shifts.NightShift)},
			}
		}

		// Validate datetime
		datetime, err := time.Parse("2006-01-02", staffType.DateTime)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"datetime": []string{"datetime should have format YYYY-MM-DD"},
			}
		}

		// Reset time to 00:00:00 UTC because we only need the date
		datetime = time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, time.UTC)
		staffType.DateTimeAsDate = &datetime

		// Validate neededStaffCount
		rxNeededStaffCountInt := regexp.MustCompile("^(?:[-+]?(?:0|[1-9][0-9]*))$")
		isValidNeededStaffCount := rxNeededStaffCountInt.MatchString(fmt.Sprintf("%v", staffType.NeededStaffCount))
		if !isValidNeededStaffCount {
			return govalidity.ValidityResponseErrors{
				"neededStaffCount": []string{"neededStaffCount can only be a valid integer number"},
			}
		}
		if staffType.NeededStaffCount < 0 {
			return govalidity.ValidityResponseErrors{
				"neededStaffCount": []string{"neededStaffCount can only be greater than 0"},
			}
		}

		// Validate roleId
		rxRoleIdInt := regexp.MustCompile("^(?:[-+]?(?:0|[1-9][0-9]*))$")
		isValidRoleId := rxRoleIdInt.MatchString(fmt.Sprintf("%v", staffType.RoleID))
		if !isValidRoleId {
			return govalidity.ValidityResponseErrors{
				"roleId": []string{"neededStaffCount can only be a valid integer number"},
			}
		}
		if staffType.RoleID < 1 {
			return govalidity.ValidityResponseErrors{
				"roleId": []string{"neededStaffCount can only be greater than 1"},
			}
		}

		// Validate startHour
		startTime, err := time.Parse("15:04", staffType.StartHour)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"startHour": []string{"startHour should have format HH:MM"},
			}
		}
		staffType.StartHourAsTime = &startTime

		// Validate endHour
		endTime, err := time.Parse("15:04", staffType.EndHour)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"endHour": []string{"endHour should have format HH:MM"},
			}
		}
		staffType.EndHourAsTime = &endTime
	}

	return nil
}
