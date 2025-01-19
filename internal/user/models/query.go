package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: UserFilterType
 */
type UserFilterType struct {
	Phone           filters.FilterValue[string] `json:"phone,omitempty" openapi:"nullable;$ref:FilterValueString;example:{\"phone\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	WorkPhoneNumber filters.FilterValue[string] `json:"workPhoneNumber,omitempty" openapi:"nullable;$ref:FilterValueString;example:{\"workPhoneNumber\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Email           filters.FilterValue[string] `json:"email,omitempty" openapi:"$ref:FilterValueString;example:{\"email\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	FirstName       filters.FilterValue[string] `json:"firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"firstName\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	LastName        filters.FilterValue[string] `json:"lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"lastName\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Username        filters.FilterValue[string] `json:"username,omitempty" openapi:"$ref:FilterValueString;example:{\"username\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	NationalCode    filters.FilterValue[string] `json:"nationalCode,omitempty" openapi:"$ref:FilterValueString;example:{\"nationalCode\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	BirthDate       filters.FilterValue[string] `json:"birthDate,omitempty" openapi:"$ref:FilterValueString;example:{\"birthDate\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	AvatarUrl       filters.FilterValue[string] `json:"avatarUrl,omitempty" openapi:"$ref:FilterValueString;example:{\"avatarUrl\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	SuspendedAt     filters.FilterValue[string] `json:"suspended_at,omitempty" openapi:"$ref:FilterValueString;example:{\"suspended_at\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt       filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"hohohoohoho\"}"`
}

/*
 * @apiDefine: UserSortValue
 */
type UserSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: UserSortType
 */
type UserSortType struct {
	FirstName    UserSortValue `json:"firstName,omitempty" openapi:"$ref:UserSortValue;example:{\"firstName\":{\"op\":\"asc\"}}"`
	LastName     UserSortValue `json:"lastName,omitempty" openapi:"$ref:UserSortValue;example:{\"lastName\":{\"op\":\"asc\"}}"`
	Username     UserSortValue `json:"username,omitempty" openapi:"$ref:UserSortValue;example:{\"username\":{\"op\":\"asc\"}}"`
	NationalCode UserSortValue `json:"nationalCode,omitempty" openapi:"$ref:UserSortValue;example:{\"nationalCode\":{\"op\":\"asc\"}}"`
	BirthDate    UserSortValue `json:"birthDate,omitempty" openapi:"$ref:UserSortValue;example:{\"birthDate\":{\"op\":\"asc\"}}"`
	AvatarUrl    UserSortValue `json:"avatarUrl,omitempty" openapi:"$ref:UserSortValue;example:{\"avatarUrl\":{\"op\":\"asc\"}}"`
	SuspendedAt  UserSortValue `json:"suspended_at,omitempty" openapi:"$ref:UserSortValue;example:{\"suspended_at\":{\"op\":\"asc\"}}"`
	CreatedAt    UserSortValue `json:"created_at,omitempty" openapi:"$ref:UserSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: UsersQueryRequestParams
 */
type UsersQueryRequestParams struct {
	ID      int            `json:"id,string,omitempty" openapi:"example:1"`
	Page    int            `json:"page,string,omitempty" openapi:"example:1;nullable;pattern:^[0-9]+$;"`
	Limit   int            `json:"limit,string,omitempty" openapi:"example:10"`
	Filters UserFilterType `json:"filters,omitempty" openapi:"$ref:UserFilterType;in:query"`
	Sorts   UserSortType   `json:"sorts,omitempty" openapi:"$ref:UserSortType;in:query"`
}

func (data *UsersQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"phone": govalidity.Schema{
				"op":    govalidity.New("filter.phone.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.phone.value").Optional(),
			},
			"workPhoneNumber": govalidity.Schema{
				"op":    govalidity.New("filter.workPhoneNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.workPhoneNumber.value").Optional(),
			},
			"email": govalidity.Schema{
				"op":    govalidity.New("filter.email.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.email.value").Email().Optional(),
			},
			"firstName": govalidity.Schema{
				"op":    govalidity.New("filter.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.firstName.value").Optional(),
			},
			"lastName": govalidity.Schema{
				"op":    govalidity.New("filter.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.lastName.value").Optional(),
			},
			"username": govalidity.Schema{
				"op":    govalidity.New("filter.username.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.username.value").Optional(),
			},
			"nationalCode": govalidity.Schema{
				"op":    govalidity.New("filter.nationalCode.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.nationalCode.value").Optional(),
			},
			"birthDate": govalidity.Schema{
				"op":    govalidity.New("filter.birthdate.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.birthdate.value").Optional(),
			},
			"avatarUrl": govalidity.Schema{
				"op":    govalidity.New("filter.avatarUrl.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.avatarUrl.value").Optional(),
			},
			"suspended_at": govalidity.Schema{
				"op":    govalidity.New("filter.suspended_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.suspended_at.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"firstName": govalidity.Schema{
				"op": govalidity.New("sort.firstName.op").In([]string{"asc", "desc"}).Optional(),
			},
			"lastName": govalidity.Schema{
				"op": govalidity.New("sort.lastName.op").In([]string{"asc", "desc"}).Optional(),
			},
			"username": govalidity.Schema{
				"op": govalidity.New("sort.username.op").In([]string{"asc", "desc"}).Optional(),
			},
			"nationalCode": govalidity.Schema{
				"op": govalidity.New("sort.nationalCode.op").In([]string{"asc", "desc"}).Optional(),
			},
			"birthDate": govalidity.Schema{
				"op": govalidity.New("sort.birthDate.op").In([]string{"asc", "desc"}).Optional(),
			},
			"avatarUrl": govalidity.Schema{
				"op": govalidity.New("sort.avatarUrl.op").In([]string{"asc", "desc"}).Optional(),
			},
			"suspended_at": govalidity.Schema{
				"op": govalidity.New("sort.suspended_at.op").In([]string{"asc", "desc"}).Optional(),
			},
			"created_at": govalidity.Schema{
				"op": govalidity.New("sort.created_at.op").In([]string{"asc", "desc"}).Optional(),
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
