package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: StaffFilterType
 */
type StaffFilterType struct {
	UserId      filters.FilterValue[int]    `json:"userId,omitempty" openapi:"$ref:FilterValueInt;example:{\"userId\":{\"op\":\"equals\",\"value\":1}"`
	FirstName   filters.FilterValue[string] `json:"user.firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"user.firstName\":{\"op\":\"equals\",\"value\":1}"`
	LastName    filters.FilterValue[string] `json:"user.lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"user.lastName\":{\"op\":\"equals\",\"value\":1}"`
	Role        filters.FilterValue[string] `json:"user.role,omitempty" openapi:"$ref:FilterValueString;example:{\"user.role\":{\"op\":\"equals\",\"value\":1}"`
	PhoneNumber filters.FilterValue[string] `json:"user.phone,omitempty" openapi:"$ref:FilterValueString;example:{\"user.phone\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Abilities   filters.FilterValue[string] `json:"abilities,omitempty" openapi:"$ref:FilterValueString;example:{\"abilities\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Team        filters.FilterValue[string] `json:"team,omitempty" openapi:"$ref:FilterValueString;example:{\"team\":{\"op\":\"equals\",\"value\":1}"`
	Grace       filters.FilterValue[string] `json:"grace,omitempty" openapi:"$ref:FilterValueString;example:{\"grace\":{\"op\":\"equals\",\"value\":1}"`
	Warning     filters.FilterValue[string] `json:"warning,omitempty" openapi:"$ref:FilterValueString;example:{\"warning\":{\"op\":\"equals\",\"value\":1}"`
	Progress    filters.FilterValue[string] `json:"progress,omitempty" openapi:"$ref:FilterValueString;example:{\"progress\":{\"op\":\"equals\",\"value\":1}"`
	CreatedAt   filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: StaffSortValue
 */
type StaffSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: StaffSortType
 */
type StaffSortType struct {
	FirstName   StaffSortValue `json:"user.firstName,omitempty" openapi:"$ref:StaffSortValue;example:{\"user.firstName\":{\"op\":\"asc\"}}"`
	LastName    StaffSortValue `json:"user.lastName,omitempty" openapi:"$ref:StaffSortValue;example:{\"user.lastName\":{\"op\":\"asc\"}}"`
	Username    StaffSortValue `json:"user.username,omitempty" openapi:"$ref:StaffSortValue;example:{\"user.username\":{\"op\":\"asc\"}}"`
	Role        StaffSortValue `json:"user.role,omitempty" openapi:"$ref:StaffSortValue;example:{\"user.role\":{\"op\":\"asc\"}}"`
	PhoneNumber StaffSortValue `json:"user.phone,omitempty" openapi:"$ref:StaffSortValue;example:{\"user.phone\":{\"op\":\"asc\"}}"`
	Abilities   StaffSortValue `json:"abilities,omitempty" openapi:"$ref:StaffSortValue;example:{\"abilities\":{\"op\":\"asc\"}}"`
	Team        StaffSortValue `json:"team,omitempty" openapi:"$ref:StaffSortValue;example:{\"team\":{\"op\":\"asc\"}}"`
	Grace       StaffSortValue `json:"grace,omitempty" openapi:"$ref:StaffSortValue;example:{\"grace\":{\"op\":\"asc\"}}"`
	Warning     StaffSortValue `json:"warning,omitempty" openapi:"$ref:StaffSortValue;example:{\"warning\":{\"op\":\"asc\"}}"`
	Progress    StaffSortValue `json:"progress,omitempty" openapi:"$ref:StaffSortValue;example:{\"progress\":{\"op\":\"asc\"}}"`
	CreatedAt   StaffSortValue `json:"created_at,omitempty" openapi:"$ref:StaffSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: StaffsQueryRequestParams
 */
type StaffsQueryRequestParams struct {
	ID      int             `json:"id,string,omitempty" openapi:"example:1"`
	UserID  int             `json:"userId,string,omitempty" openapi:"example:1"`
	Page    int             `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int             `json:"limit,string,omitempty" openapi:"example:10"`
	Filters StaffFilterType `json:"filters,omitempty" openapi:"$ref:StaffFilterType;in:query"`
	Sorts   StaffSortType   `json:"sorts,omitempty" openapi:"$ref:StaffSortType;in:query"`
}

func (data *StaffsQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"user.firstName": govalidity.Schema{
				"op":    govalidity.New("filter.user.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.firstName.value").Optional(),
			},
			"user.lastName": govalidity.Schema{
				"op":    govalidity.New("filter.user.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.lastName.value").Optional(),
			},
			"user.role": govalidity.Schema{
				"op":    govalidity.New("filter.user.role.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.role.value").Optional(),
			},
			"user.phone": govalidity.Schema{
				"op":    govalidity.New("filter.user.phone.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user.phone.value").Optional(),
			},
			"abilities": govalidity.Schema{
				"op":    govalidity.New("filter.abilities.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.abilities.value").Optional(),
			},
			"team": govalidity.Schema{
				"op":    govalidity.New("filter.team.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.team.value").Optional(),
			},
			"grace": govalidity.Schema{
				"op":    govalidity.New("filter.grace.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.grace.value").Optional(),
			},
			"warning": govalidity.Schema{
				"op":    govalidity.New("filter.warning.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.warning.value").Optional(),
			},
			"progress": govalidity.Schema{
				"op":    govalidity.New("filter.progress.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.progress.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"user.firstName": govalidity.Schema{
				"op": govalidity.New("sort.user.firstName.op"),
			},
			"user.lastName": govalidity.Schema{
				"op": govalidity.New("sort.user.lastName.op"),
			},
			"user.username": govalidity.Schema{
				"op": govalidity.New("sort.user.username.op"),
			},
			"user.role": govalidity.Schema{
				"op": govalidity.New("sort.user.role.op"),
			},
			"user.phone": govalidity.Schema{
				"op": govalidity.New("sort.user.phone.op"),
			},
			"created_at": govalidity.Schema{
				"op": govalidity.New("sort.created_at.op"),
			},
		},
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
