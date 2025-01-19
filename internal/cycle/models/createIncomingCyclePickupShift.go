package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateIncomingCyclePickupShiftRequestBody
 */
type CyclesCreateIncomingCyclePickupShiftRequestBody struct {
	CycleID                    int                                                     `json:"cycleId" openapi:"example:1;required;"`
	StaffID                    int                                                     `json:"staffId" openapi:"example:1;required;"`
	CycleNextStaffTypeIDs      interface{}                                             `json:"cycleNextStaffTypeIds" openapi:"example:[1,2,3];type:array;required;"`
	DateTime                   string                                                  `json:"datetime" openapi:"example:2021-01-01;required;"`
	Staff                      *domain.CycleIncomingCyclePickupShiftStaff              `json:"-" openapi:"ignored"`
	CycleNextStaffType         *domain.CycleIncomingCyclePickupShiftCycleNextStaffType `json:"-" openapi:"ignored"`
	CycleNextStaffTypeIDsInt64 []int64                                                 `json:"-" openapi:"ignored"`
	DateTimeAsDate             *time.Time                                              `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of the incoming HTTP request.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesCreateIncomingCyclePickupShiftRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId":               govalidity.New("cycleId").Int().Min(1).Required(),
		"staffId":               govalidity.New("staffId").Int().Min(1).Required(),
		"cycleNextStaffTypeIds": govalidity.New("cycleNextStaffTypeIds"),
		"datetime":              govalidity.New("datetime").Required(),
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
	cycleNextStaffTypeIDsInt64, err := utils.ConvertInterfaceSliceToSliceOfInt64(data.CycleNextStaffTypeIDs)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"cycleStaffTypeIds": []string{"CycleStaffTypeIDs is invalid"},
		}
	}
	data.CycleNextStaffTypeIDsInt64 = cycleNextStaffTypeIDsInt64

	return nil
}
