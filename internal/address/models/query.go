package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"

	"github.com/hoitek/Go-Quilder/filters"
)

/*
 * @apiDefine: AddressFilterType
 */
type AddressFilterType struct {
	Name           filters.FilterValue[string] `json:"name,omitempty" openapi:"$ref:FilterValueString;example:{\"name\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	Street         filters.FilterValue[string] `json:"street,omitempty" openapi:"$ref:FilterValueString;example:{\"street\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	PostalCode     filters.FilterValue[string] `json:"postalCode,omitempty" openapi:"$ref:FilterValueString;example:{\"postalCode\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	BuildingNumber filters.FilterValue[string] `json:"buildingNumber,omitempty" openapi:"$ref:FilterValueString;example:{\"buildingNumber\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
	CreatedAt      filters.FilterValue[string] `json:"createdAt,omitempty" openapi:"$ref:FilterValueString;example:{\"createdAt\":{\"op\":\"equals\",\"value\":\"09123456789\"}"`
}

/*
 * @apiDefine: AddressSortValue
 */
type AddressSortValue struct {
	Op string `json:"op,omitempty" openapi:"example:asc"`
}

/*
 * @apiDefine: AddressSortType
 */
type AddressSortType struct {
	Name           AddressSortValue `json:"name,omitempty" openapi:"$ref:AddressSortValue;example:{\"name\":{\"op\":\"asc\"}}"`
	Street         AddressSortValue `json:"street,omitempty" openapi:"$ref:AddressSortValue;example:{\"street\":{\"op\":\"asc\"}}"`
	PostalCode     AddressSortValue `json:"postalCode,omitempty" openapi:"$ref:AddressSortValue;example:{\"postalCode\":{\"op\":\"asc\"}}"`
	BuildingNumber AddressSortValue `json:"buildingNumber,omitempty" openapi:"$ref:AddressSortValue;example:{\"buildingNumber\":{\"op\":\"asc\"}}"`
	CreatedAt      AddressSortValue `json:"created_at,omitempty" openapi:"$ref:AddressSortValue;example:{\"created_at\":{\"op\":\"asc\"}}"`
}

/*
 * @apiDefine: AddressesQueryRequestParams
 */
type AddressesQueryRequestParams struct {
	ID         int               `json:"id,string,omitempty" openapi:"example:1"`
	StaffID    int               `json:"staffId,string,omitempty" openapi:"example:1"`
	CustomerID int               `json:"customerId,string,omitempty" openapi:"example:1"`
	Page       int               `json:"page,string,omitempty" openapi:"example:1"`
	Limit      int               `json:"limit,string,omitempty" openapi:"example:10"`
	Filters    AddressFilterType `json:"filters,omitempty" openapi:"$ref:AddressFilterType;in:query"`
	Sorts      AddressSortType   `json:"sorts,omitempty" openapi:"$ref:AddressSortType;in:query"`
}

func (data *AddressesQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Optional(),
		"page":       govalidity.New("page").Int().Default("1"),
		"limit":      govalidity.New("limit").Int().Default("10"),
		"staffId":    govalidity.New("staffId").Int().Optional(),
		"customerId": govalidity.New("customerId").Int().Optional(),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"street": govalidity.Schema{
				"op":    govalidity.New("filter.street.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.street.value").Optional(),
			},
			"postalCode": govalidity.Schema{
				"op":    govalidity.New("filter.postalCode.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.postalCode.value").Optional(),
			},
			"buildingNumber": govalidity.Schema{
				"op":    govalidity.New("filter.buildingNumber.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.buildingNumber.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
		},
		"sorts": govalidity.Schema{
			"name": govalidity.Schema{
				"op": govalidity.New("sort.name.op"),
			},
			"street": govalidity.Schema{
				"op": govalidity.New("sort.street.op"),
			},
			"postalCode": govalidity.Schema{
				"op": govalidity.New("sort.postalCode.op"),
			},
			"buildingNumber": govalidity.Schema{
				"op": govalidity.New("sort.buildingNumber.op"),
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
