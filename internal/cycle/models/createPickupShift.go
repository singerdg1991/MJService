package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreatePickupShiftRequestBody
 */
type CyclesCreatePickupShiftRequestBody struct {
	CycleID                int                           `json:"cycleId" openapi:"example:1;required;"`
	StaffID                int                           `json:"staffId" openapi:"example:1;required;"`
	CycleStaffTypeIDs      interface{}                   `json:"cycleStaffTypeIds" openapi:"example:[1,2,3];type:array;required;"`
	DateTime               string                        `json:"datetime" openapi:"example:2021-01-01;required;"`
	Staff                  *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
	CycleStaffTypeIDsInt64 []int64                       `json:"-" openapi:"ignored"`
	DateTimeAsDate         *time.Time                    `json:"-" openapi:"ignored"`
	ShiftName              string                        `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreatePickupShiftRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId":           govalidity.New("cycleId").Int().Min(1).Required(),
		"staffId":           govalidity.New("staffId").Int().Min(1).Required(),
		"cycleStaffTypeIds": govalidity.New("cycleStaffTypeIds"),
		"datetime":          govalidity.New("datetime").Required(),
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

	// Validate cycleStaffTypeIDs
	cycleStaffTypeIDsInt64, err := utils.ConvertInterfaceSliceToSliceOfInt64(data.CycleStaffTypeIDs)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"cycleStaffTypeIds": []string{"CycleStaffTypeIDs is invalid"},
		}
	}
	data.CycleStaffTypeIDsInt64 = cycleStaffTypeIDsInt64

	return nil
}
