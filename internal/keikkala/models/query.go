package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: KeikkalaFilterType
 */
type KeikkalaFilterType struct {
	Title           filters.FilterValue[string] `json:"title,omitempty" openapi:"$ref:FilterValueString;example:{\"title\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	SectionIDs      filters.FilterValue[string] `json:"sectionIds,omitempty" openapi:"$ref:FilterValueString;example:{\"sectionIds\":{\"op\":\"equals\",\"value\":\"[1,2,3]\"}"`
	RoleID          filters.FilterValue[int]    `json:"roleId,omitempty" openapi:"$ref:FilterValueInt;example:{\"roleId\":{\"op\":\"equals\",\"value\":\"1\"}"`
	RoleName        filters.FilterValue[string] `json:"roleName,omitempty" openapi:"$ref:FilterValueString;example:{\"roleName\":{\"op\":\"equals\",\"value\":\"1\"}"`
	StartDate       filters.FilterValue[string] `json:"start_date,omitempty" openapi:"$ref:FilterValueString;example:{\"start_date\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	EndDate         filters.FilterValue[string] `json:"end_date,omitempty" openapi:"$ref:FilterValueString;example:{\"end_date\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	StartTime       filters.FilterValue[string] `json:"start_time,omitempty" openapi:"$ref:FilterValueString;example:{\"start_time\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	EndTime         filters.FilterValue[string] `json:"end_time,omitempty" openapi:"$ref:FilterValueString;example:{\"end_time\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	KaupunkiAddress filters.FilterValue[string] `json:"kaupunkiAddress,omitempty" openapi:"$ref:FilterValueString;example:{\"kaupunkiAddress\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	PaymentType     filters.FilterValue[string] `json:"paymentType,omitempty" openapi:"$ref:FilterValueString;example:{\"paymentType\":{\"op\":\"equals\",\"value\":\"paySoon\"}"`
	ShiftName       filters.FilterValue[string] `json:"shiftName,omitempty" openapi:"$ref:FilterValueString;example:{\"shiftName\":{\"op\":\"equals\",\"value\":\"morning\"}"`
	Description     filters.FilterValue[string] `json:"description,omitempty" openapi:"$ref:FilterValueString;example:{\"description\":{\"op\":\"equals\",\"value\":\"description\"}"`
	Status          filters.FilterValue[string] `json:"status,omitempty" openapi:"$ref:FilterValueString;example:{\"status\":{\"op\":\"equals\",\"value\":\"open\"}"`
	PickedAt        filters.FilterValue[string] `json:"picked_at,omitempty" openapi:"$ref:FilterValueString;example:{\"picked_at\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	PickedBy        filters.FilterValue[int]    `json:"pickedBy,omitempty" openapi:"$ref:FilterValueInt;example:{\"pickedBy\":{\"op\":\"equals\",\"value\":\"1\"}"`
	CreatedAt       filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
	UpdatedAt       filters.FilterValue[string] `json:"updated_at,omitempty" openapi:"$ref:FilterValueString;example:{\"updated_at\":{\"op\":\"equals\",\"value\":\"2021-01-01\"}"`
}

/*
 * @apiDefine: KeikkalaSortValue
 */
type KeikkalaSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: KeikkalaSortType
 */
type KeikkalaSortType struct {
	ID        KeikkalaSortValue `json:"id,omitempty" openapi:"$ref:KeikkalaSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	StartDate KeikkalaSortValue `json:"start_date,omitempty" openapi:"$ref:KeikkalaSortValue;example:{\"start_date\":{\"op\":\"asc\"}}"`
	CreatedAt KeikkalaSortValue `json:"created_at,omitempty" openapi:"$ref:KeikkalaSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: KeikkalasQueryRequestParams
 */
type KeikkalasQueryRequestParams struct {
	ID      int                `json:"id,string,omitempty" openapi:"example:1"`
	Page    int                `json:"page,string,omitempty" openapi:"example:1"`
	Limit   int                `json:"limit,string,omitempty" openapi:"example:10"`
	Filters KeikkalaFilterType `json:"filters,omitempty" openapi:"$ref:KeikkalaFilterType;in:query"`
	Sorts   KeikkalaSortType   `json:"sorts,omitempty" openapi:"$ref:KeikkalaSortType;in:query"`
}

func (data *KeikkalasQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":    govalidity.New("id").Int().Optional(),
		"page":  govalidity.New("page").Int().Default("1"),
		"limit": govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"sectionIds": govalidity.Schema{
				"op":    govalidity.New("filter.sectionIds.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.sectionIds.value").Optional(),
			},
			"roleId": govalidity.Schema{
				"op":    govalidity.New("filter.roleId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.roleId.value").Optional(),
			},
			"roleName": govalidity.Schema{
				"op":    govalidity.New("filter.roleName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.roleName.value").Optional(),
			},
			"start_date": govalidity.Schema{
				"op":    govalidity.New("filter.start_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_date.value").Optional(),
			},
			"end_date": govalidity.Schema{
				"op":    govalidity.New("filter.end_date.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.end_date.value").Optional(),
			},
			"start_time": govalidity.Schema{
				"op":    govalidity.New("filter.start_time.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.start_time.value").Optional(),
			},
			"end_time": govalidity.Schema{
				"op":    govalidity.New("filter.end_time.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.end_time.value").Optional(),
			},
			"kaupunkiAddress": govalidity.Schema{
				"op":    govalidity.New("filter.kaupunkiAddress.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.kaupunkiAddress.value").Optional(),
			},
			"paymentType": govalidity.Schema{
				"op":    govalidity.New("filter.paymentType.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.paymentType.value").Optional(),
			},
			"description": govalidity.Schema{
				"op":    govalidity.New("filter.description.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.description.value").Optional(),
			},
			"status": govalidity.Schema{
				"op":    govalidity.New("filter.status.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.status.value").Optional(),
			},
			"picked_at": govalidity.Schema{
				"op":    govalidity.New("filter.picked_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.picked_at.value").Optional(),
			},
			"picked_by": govalidity.Schema{
				"op":    govalidity.New("filter.picked_by.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.picked_by.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"updated_at": govalidity.Schema{
				"op":    govalidity.New("filter.updated_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.updated_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"start_date": govalidity.Schema{
				"op": govalidity.New("sort.start_date.op"),
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
