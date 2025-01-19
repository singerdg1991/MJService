package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsFilterType
 */
type CustomersQueryContractualMobilityRestrictionLogsFilterType struct {
	BeforeValue filters.FilterValue[string] `json:"beforeValue,omitempty" openapi:"$ref:FilterValueString;example:{\"beforeValue\":{\"op\":\"equals\",\"value\":\"BeforeValue\"}}"`
	AfterValue  filters.FilterValue[string] `json:"afterValue,omitempty" openapi:"$ref:FilterValueString;example:{\"afterValue\":{\"op\":\"equals\",\"value\":\"AfterValue\"}}"`
	FirstName   filters.FilterValue[string] `json:"firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"firstName\":{\"op\":\"equals\",\"value\":\"firstName\"}}"`
	LastName    filters.FilterValue[string] `json:"lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"lastName\":{\"op\":\"equals\",\"value\":\"lastName\"}}"`
	FullName    filters.FilterValue[string] `json:"fullName,omitempty" openapi:"$ref:FilterValueString;example:{\"fullName\":{\"op\":\"equals\",\"value\":\"fullName\"}}"`
	CreatedAt   filters.FilterValue[string] `json:"created_at,omitempty" openapi:"$ref:FilterValueString;example:{\"created_at\":{\"op\":\"equals\",\"value\":\"2021-01-01T00:00:00Z\"}}"`
}

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsSortValue
 */
type CustomersQueryContractualMobilityRestrictionLogsSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsSortType
 */
type CustomersQueryContractualMobilityRestrictionLogsSortType struct {
	ID          CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"id,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	BeforeValue CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"beforeValue,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"beforeValue\":{\"op\":\"asc\"}}"`
	AfterValue  CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"afterValue,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"afterValue\":{\"op\":\"asc\"}}"`
	FirstName   CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"firstName,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"firstName\":{\"op\":\"asc\"}}"`
	LastName    CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"lastName,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"lastName\":{\"op\":\"asc\"}}"`
	CreatedAt   CustomersQueryContractualMobilityRestrictionLogsSortValue `json:"created_at,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryContractualMobilityRestrictionLogsRequestParams
 */
type CustomersQueryContractualMobilityRestrictionLogsRequestParams struct {
	ID         int                                                        `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                                                        `json:"customerId,string,omitempty" openapi:"example:1;required"`
	Page       int                                                        `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                                                        `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryContractualMobilityRestrictionLogsFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsFilterType;in:query"`
	Sorts      CustomersQueryContractualMobilityRestrictionLogsSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryContractualMobilityRestrictionLogsSortType;in:query"`
}

func (data *CustomersQueryContractualMobilityRestrictionLogsRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Required(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"beforeValue": govalidity.Schema{
				"op":    govalidity.New("filter.beforeValue.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.beforeValue.value").Optional(),
			},
			"afterValue": govalidity.Schema{
				"op":    govalidity.New("filter.afterValue.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.afterValue.value").Optional(),
			},
			"firstName": govalidity.Schema{
				"op":    govalidity.New("filter.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.firstName.value").Optional(),
			},
			"lastName": govalidity.Schema{
				"op":    govalidity.New("filter.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.lastName.value").Optional(),
			},
			"fullName": govalidity.Schema{
				"op":    govalidity.New("filter.fullName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.fullName.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"beforeValue": govalidity.Schema{
				"op": govalidity.New("sort.beforeValue.op"),
			},
			"afterValue": govalidity.Schema{
				"op": govalidity.New("sort.afterValue.op"),
			},
			"firstName": govalidity.Schema{
				"op": govalidity.New("sort.firstName.op"),
			},
			"lastName": govalidity.Schema{
				"op": govalidity.New("sort.lastName.op"),
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
