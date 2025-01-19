package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: CustomersQueryRelativesFilterType
 */
type CustomersQueryRelativesFilterType struct {
	CustomerID            filters.FilterValue[int]    `json:"customerId,omitempty" openapi:"$ref:FilterValueInt;example:{\"customerId\":{\"op\":\"equals\",\"value\":1}}"`
	CityID                filters.FilterValue[int]    `json:"cityId,omitempty" openapi:"$ref:FilterValueInt;example:{\"cityId\":{\"op\":\"equals\",\"value\":1}}"`
	FirstName             filters.FilterValue[string] `json:"firstName,omitempty" openapi:"$ref:FilterValueString;example:{\"firstName\":{\"op\":\"equals\",\"value\":\"John\"}}"`
	LastName              filters.FilterValue[string] `json:"lastName,omitempty" openapi:"$ref:FilterValueString;example:{\"lastName\":{\"op\":\"equals\",\"value\":\"Doe\"}}"`
	PhoneNumber           filters.FilterValue[string] `json:"phoneNumber,omitempty" openapi:"$ref:FilterValueString;example:{\"phoneNumber\":{\"op\":\"equals\",\"value\":\"09123456789\"}}"`
	Relation              filters.FilterValue[string] `json:"relation,omitempty" openapi:"$ref:FilterValueString;example:{\"relation\":{\"op\":\"equals\",\"value\":\"father\"}}"`
	AddressName           filters.FilterValue[string] `json:"addressName,omitempty" openapi:"$ref:FilterValueString;example:{\"addressName\":{\"op\":\"equals\",\"value\":\"home\"}}"`
	AddressStreet         filters.FilterValue[string] `json:"addressStreet,omitempty" openapi:"$ref:FilterValueString;example:{\"addressStreet\":{\"op\":\"equals\",\"value\":\"street\"}}"`
	AddressBuildingNumber filters.FilterValue[string] `json:"addressBuildingNumber,omitempty" openapi:"$ref:FilterValueString;example:{\"addressBuildingNumber\":{\"op\":\"equals\",\"value\":\"1\"}}"`
	AddressPostalCode     filters.FilterValue[string] `json:"addressPostalCode,omitempty" openapi:"$ref:FilterValueString;example:{\"addressPostalCode\":{\"op\":\"equals\",\"value\":\"1234567890\"}}"`
}

/*
 * @apiDefine: CustomersQueryRelativesSortValue
 */
type CustomersQueryRelativesSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: CustomersQueryRelativesSortType
 */
type CustomersQueryRelativesSortType struct {
	ID                    CustomersQueryRelativesSortValue `json:"id,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"id\":{\"op\":\"asc\"}}"`
	FirstName             CustomersQueryRelativesSortValue `json:"firstName,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"firstName\":{\"op\":\"asc\"}}"`
	LastName              CustomersQueryRelativesSortValue `json:"lastName,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"lastName\":{\"op\":\"asc\"}}"`
	PhoneNumber           CustomersQueryRelativesSortValue `json:"phoneNumber,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"phoneNumber\":{\"op\":\"asc\"}}"`
	Relation              CustomersQueryRelativesSortValue `json:"relation,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"relation\":{\"op\":\"asc\"}}"`
	AddressName           CustomersQueryRelativesSortValue `json:"addressName,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"addressName\":{\"op\":\"asc\"}}"`
	AddressStreet         CustomersQueryRelativesSortValue `json:"addressStreet,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"addressStreet\":{\"op\":\"asc\"}}"`
	AddressBuildingNumber CustomersQueryRelativesSortValue `json:"addressBuildingNumber,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"addressBuildingNumber\":{\"op\":\"asc\"}}"`
	AddressPostalCode     CustomersQueryRelativesSortValue `json:"addressPostalCode,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"addressPostalCode\":{\"op\":\"asc\"}}"`
	CreatedAt             CustomersQueryRelativesSortValue `json:"created_at,omitempty" openapi:"$ref:CustomersQueryRelativesSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: CustomersQueryRelativesRequestParams
 */
type CustomersQueryRelativesRequestParams struct {
	ID         int                               `json:"id,string,omitempty" openapi:"example:1"`
	CustomerID int                               `json:"customerId,string,omitempty" openapi:"example:1;required"`
	Page       int                               `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int                               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    CustomersQueryRelativesFilterType `json:"filters,omitempty" openapi:"$ref:CustomersQueryRelativesFilterType;in:query"`
	Sorts      CustomersQueryRelativesSortType   `json:"sorts,omitempty" openapi:"$ref:CustomersQueryRelativesSortType;in:query"`
}

func (data *CustomersQueryRelativesRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Required(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"filters": govalidity.Schema{
			"customerId": govalidity.Schema{
				"op":    govalidity.New("filter.customerId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.customerId.value").Optional(),
			},
			"cityId": govalidity.Schema{
				"op":    govalidity.New("filter.cityId.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.cityId.value").Optional(),
			},
			"firstName": govalidity.Schema{
				"op":    govalidity.New("filter.firstName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.firstName.value").Optional(),
			},
			"lastName": govalidity.Schema{
				"op":    govalidity.New("filter.lastName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.lastName.value").Optional(),
			},
			"phoneNumber": govalidity.Schema{
				"op":    govalidity.New("filter.phoneNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.phoneNumber.value").Optional(),
			},
			"relation": govalidity.Schema{
				"op":    govalidity.New("filter.relation.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.relation.value").Optional(),
			},
			"addressName": govalidity.Schema{
				"op":    govalidity.New("filter.addressName.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.addressName.value").Optional(),
			},
			"addressStreet": govalidity.Schema{
				"op":    govalidity.New("filter.addressStreet.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.addressStreet.value").Optional(),
			},
			"addressBuildingNumber": govalidity.Schema{
				"op":    govalidity.New("filter.addressBuildingNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.addressBuildingNumber.value").Optional(),
			},
			"addressPostalCode": govalidity.Schema{
				"op":    govalidity.New("filter.addressPostalCode.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.addressPostalCode.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"id": govalidity.Schema{
				"op": govalidity.New("sort.id.op"),
			},
			"firstName": govalidity.Schema{
				"op": govalidity.New("sort.firstName.op"),
			},
			"lastName": govalidity.Schema{
				"op": govalidity.New("sort.lastName.op"),
			},
			"phoneNumber": govalidity.Schema{
				"op": govalidity.New("sort.phoneNumber.op"),
			},
			"relation": govalidity.Schema{
				"op": govalidity.New("sort.relation.op"),
			},
			"addressName": govalidity.Schema{
				"op": govalidity.New("sort.addressName.op"),
			},
			"addressStreet": govalidity.Schema{
				"op": govalidity.New("sort.addressStreet.op"),
			},
			"addressBuildingNumber": govalidity.Schema{
				"op": govalidity.New("sort.addressBuildingNumber.op"),
			},
			"addressPostalCode": govalidity.Schema{
				"op": govalidity.New("sort.addressPostalCode.op"),
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
