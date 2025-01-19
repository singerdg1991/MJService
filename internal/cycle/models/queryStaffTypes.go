package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: CycleQueryStaffTypesFilterType
 */
type CycleQueryStaffTypesFilterType struct {
	ShiftName     filters.FilterValue[string] `json:"shiftName,omitempty" openapi:"$ref:FilterValueString;example:{\"shiftName\":{\"op\":\"equals\",\"value\":\"morning\"}"`
	DateTime      filters.FilterValue[string] `json:"datetime,omitempty" openapi:"$ref:FilterValueString;example:{\"datetime\":{\"op\":\"equals\",\"value\":\"2020-01-01\"}"`
	RoleID        filters.FilterValue[int]    `json:"roleId,omitempty" openapi:"$ref:FilterValueInt;example:{\"roleId\":{\"op\":\"equals\",\"value\":\"1\"}"`
	RoleName      filters.FilterValue[string] `json:"roleName,omitempty" openapi:"$ref:FilterValueString;example:{\"roleName\":{\"op\":\"equals\",\"value\":\"doctor\"}"`
	StartHour     filters.FilterValue[string] `json:"startHour,omitempty" openapi:"$ref:FilterValueString;example:{\"startHour\":{\"op\":\"equals\",\"value\":\"00:00\"}"`
	EndHour       filters.FilterValue[string] `json:"endHour,omitempty" openapi:"$ref:FilterValueString;example:{\"endHour\":{\"op\":\"equals\",\"value\":\"00:00\"}"`
	DateRangeFrom filters.FilterValue[string] `json:"dateRangeFrom,omitempty" openapi:"$ref:FilterValueString;example:{\"dateRangeFrom\":{\"op\":\"equals\",\"value\":\"2020-01-01\"}"`
	DateRangeTo   filters.FilterValue[string] `json:"dateRangeTo,omitempty" openapi:"$ref:FilterValueString;example:{\"dateRangeTo\":{\"op\":\"equals\",\"value\":\"2020-01-01\"}"`
}

/*
 * @apiDefine: CycleQueryStaffTypesSortValue
 */
type CycleQueryStaffTypesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CycleQueryStaffTypesSortType
 */
type CycleQueryStaffTypesSortType struct {
}

/*
 * @apiDefine: CyclesQueryStaffTypesRequestParams
 */
type CyclesQueryStaffTypesRequestParams struct {
	ID           int                            `json:"id,string,omitempty" openapi:"example:1"`
	CycleID      int                            `json:"cycleid,string,omitempty" openapi:"example:1"`
	RoleIDs      interface{}                    `json:"roleids,omitempty" openapi:"example:[1,2,3]"`
	Page         int                            `json:"page,string,omitempty" openapi:"example:1"`
	Limit        int                            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters      CycleQueryStaffTypesFilterType `json:"filters,omitempty" openapi:"$ref:CycleQueryStaffTypesFilterType;in:query"`
	Sorts        CycleQueryStaffTypesSortType   `json:"sorts,omitempty" openapi:"$ref:CycleQueryStaffTypesSortType;in:query"`
	RoleIDsInt64 []int64                        `json:"-" openapi:"ignored"` // This is a helper field to convert RoleIDs to int64
}

// ValidateQueries validates the queries in the request.
//
// It takes an http.Request as a parameter and returns a map of validation errors.
// The function checks the schema of the request parameters and returns an error if any of the parameters are invalid.
// It also checks the datetime format and the RoleIDs format.
//
// Parameter: r *http.Request
// Return type: govalidity.ValidityResponseErrors
func (data *CyclesQueryStaffTypesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"cycleid": govalidity.New("id").Int().Optional(),
		"roleids": govalidity.New("roleids").Optional(),
		"page":    govalidity.New("page").Int().Default("1"),
		"limit":   govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"shiftName": govalidity.Schema{
				"op":    govalidity.New("filter.shiftName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.shiftName.value").Optional(),
			},
			"datetime": govalidity.Schema{
				"op":    govalidity.New("filter.datetime.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.datetime.value").Optional(),
			},
			"roleId": govalidity.Schema{
				"op":    govalidity.New("filter.roleId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.roleId.value").Optional(),
			},
			"roleName": govalidity.Schema{
				"op":    govalidity.New("filter.roleName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.roleName.value").Optional(),
			},
			"startHour": govalidity.Schema{
				"op":    govalidity.New("filter.startHour.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.startHour.value").Optional(),
			},
			"endHour": govalidity.Schema{
				"op":    govalidity.New("filter.endHour.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.endHour.value").Optional(),
			},
			"dateRangeFrom": govalidity.Schema{
				"op":    govalidity.New("filter.dateRangeFrom.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.dateRangeFrom.value").Optional(),
			},
			"dateRangeTo": govalidity.Schema{
				"op":    govalidity.New("filter.dateRangeTo.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.dateRangeTo.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{},
	}

	// Validate queries
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

	// Validate RoleIDs
	if data.RoleIDs != nil {
		err := json.Unmarshal([]byte(data.RoleIDs.(string)), &data.RoleIDsInt64)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"roleids": []string{"roleids should be an array of integers"},
			}
		}
	}

	return nil
}
