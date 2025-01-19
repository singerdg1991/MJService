package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CycleQueryIncomingCyclePickupShiftsFilterType
 */
type CycleQueryIncomingCyclePickupShiftsFilterType struct {
	ShiftName filters.FilterValue[string] `json:"shiftName,omitempty" openapi:"$ref:FilterValueString;example:{\"shiftName\":{\"op\":\"equals\",\"value\":\"morning\"}"`
	DateTime  filters.FilterValue[string] `json:"datetime,omitempty" openapi:"$ref:FilterValueString;example:{\"datetime\":{\"op\":\"equals\",\"value\":\"2020-01-01\"}"`
	StartHour filters.FilterValue[string] `json:"startHour,omitempty" openapi:"$ref:FilterValueString;example:{\"startHour\":{\"op\":\"equals\",\"value\":\"00:00\"}"`
	EndHour   filters.FilterValue[string] `json:"endHour,omitempty" openapi:"$ref:FilterValueString;example:{\"endHour\":{\"op\":\"equals\",\"value\":\"00:00\"}"`
}

/*
 * @apiDefine: CycleQueryIncomingCyclePickupShiftsSortValue
 */
type CycleQueryIncomingCyclePickupShiftsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CycleQueryIncomingCyclePickupShiftsSortType
 */
type CycleQueryIncomingCyclePickupShiftsSortType struct {
}

/*
 * @apiDefine: CyclesQueryIncomingCyclePickupShiftsRequestParams
 */
type CyclesQueryIncomingCyclePickupShiftsRequestParams struct {
	ID                         int                                           `json:"id,string,omitempty" openapi:"example:1"`
	CycleID                    int                                           `json:"cycleid,string,omitempty" openapi:"example:1"`
	StaffID                    int                                           `json:"staffid,string,omitempty" openapi:"example:1"`
	CycleNextStaffTypeIDs      string                                        `json:"cyclenextstafftypeids,omitempty" openapi:"example:[1,2,3]"`
	ShiftNames                 string                                        `json:"shiftnames,omitempty" openapi:"example:[\"morning\",\"evening\",\"night\"]"`
	RangeDateTimeStart         string                                        `json:"rangedatetimestart,omitempty" openapi:"example:2020-01-01"`
	RangeDateTimeEnd           string                                        `json:"rangedatetimeend,omitempty" openapi:"example:2020-01-01"`
	Page                       int                                           `json:"page,string,omitempty" openapi:"example:1"`
	Limit                      int                                           `json:"limit,string,omitempty" openapi:"example:10"`
	Filters                    CycleQueryIncomingCyclePickupShiftsFilterType `json:"filters,omitempty" openapi:"$ref:CycleQueryIncomingCyclePickupShiftsFilterType;in:query"`
	Sorts                      CycleQueryIncomingCyclePickupShiftsSortType   `json:"sorts,omitempty" openapi:"$ref:CycleQueryIncomingCyclePickupShiftsSortType;in:query"`
	CycleNextStaffTypeIDsInt64 []int64                                       `json:"-" openapi:"ignored"`
	ShiftNamesAsArray          []string                                      `json:"-" openapi:"ignored"`
	RangeDateTimeStartAsTime   *time.Time                                    `json:"-" openapi:"ignored"`
	RangeDateTimeEndAsTime     *time.Time                                    `json:"-" openapi:"ignored"`
}

// ValidateQueries validates the queries of the CyclesQueryIncomingCyclePickupShiftsRequestParams.
//
// It takes an http.Request object as a parameter.
// Returns a govalidity.ValidityResponseErrors object containing any validation errors.
func (data *CyclesQueryIncomingCyclePickupShiftsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":                    govalidity.New("id").Int().Optional(),
		"cycleid":               govalidity.New("cycleid").Int().Optional(),
		"staffid":               govalidity.New("staffid").Int().Optional(),
		"cyclenextstafftypeids": govalidity.New("cyclenextstafftypeids").Optional(),
		"rangedatetimestart":    govalidity.New("rangedatetimestart").Optional(),
		"rangedatetimeend":      govalidity.New("rangedatetimeend").Optional(),
		"shiftnames":            govalidity.New("shiftnames").Optional(),
		"page":                  govalidity.New("page").Int().Default("1"),
		"limit":                 govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"shiftName": govalidity.Schema{
				"op":    govalidity.New("filter.shiftName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.shiftName.value").Optional(),
			},
			"datetime": govalidity.Schema{
				"op":    govalidity.New("filter.datetime.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.datetime.value").Optional(),
			},
			"staffId": govalidity.Schema{
				"op":    govalidity.New("filter.staffId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.staffId.value").Optional(),
			},
			"startHour": govalidity.Schema{
				"op":    govalidity.New("filter.startHour.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.startHour.value").Optional(),
			},
			"endHour": govalidity.Schema{
				"op":    govalidity.New("filter.endHour.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.endHour.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate datetime
	datetime := data.Filters.DateTime.Value
	if datetime != "" {
		_, err := time.Parse("2006-01-02", datetime)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"filter.datetime.value": []string{"datetime should have format YYYY-MM-DD"},
			}
		}
	}

	// Validate CycleNextStaffTypeIDs if it exists
	if data.CycleNextStaffTypeIDs != "" {
		err := json.Unmarshal([]byte(data.CycleNextStaffTypeIDs), &data.CycleNextStaffTypeIDsInt64)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"cyclenextstafftypeids": []string{"CycleNextStaffTypeIDs should be a valid JSON array of integers"},
			}
		}
	}

	// Validate ShiftNames if it exists
	if data.ShiftNames != "" {
		err := json.Unmarshal([]byte(data.ShiftNames), &data.ShiftNamesAsArray)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"shiftNames": []string{"ShiftNames should be a valid JSON array of strings"},
			}
		}
	}

	// Validate RangeDateTimeStart if it exists
	if data.RangeDateTimeStart != "" {
		ts, err := time.Parse("2006-01-02", data.RangeDateTimeStart)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"rangedatetimestart": []string{"rangedatetimestart should have format YYYY-MM-DD"},
			}
		}
		data.RangeDateTimeStartAsTime = &ts
	}

	// Validate RangeDateTimeEnd if it exists
	if data.RangeDateTimeEnd != "" {
		te, err := time.Parse("2006-01-02", data.RangeDateTimeEnd)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"rangedatetimeend": []string{"rangedatetimeend should have format YYYY-MM-DD"},
			}
		}
		data.RangeDateTimeEndAsTime = &te
	}

	return nil
}
