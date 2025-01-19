package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: AddressesCreateRequestBodyStreet
 */
type AddressesCreateRequestBodyStreet struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:tehran"`
}

/*
 * @apiDefine: AddressesCreateRequestBodyCity
 */
type AddressesCreateRequestBodyCity struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:tehran"`
}

/*
 * @apiDefine: AddressesCreateRequestBody
 */
type AddressesCreateRequestBody struct {
	StaffID                 *int                           `json:"staffId" openapi:"example:1;required;"`
	CustomerID              *int                           `json:"customerId" openapi:"example:1;required;"`
	City                    AddressesCreateRequestBodyCity `json:"city" openapi:"$ref:AddressesCreateRequestBodyCity"`
	Street                  string                         `json:"street" openapi:"example:tehran;"`
	Name                    string                         `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	PostalCode              string                         `json:"postalCode" openapi:"example:1234567890;"`
	BuildingNumber          string                         `json:"buildingNumber" openapi:"example:123;required;"`
	IsDeliveryAddress       string                         `json:"isDeliveryAddress" openapi:"example:true;required;"`
	IsMainAddress           string                         `json:"isMainAddress" openapi:"example:true;required;"`
	IsDeliveryAddressAsBool bool                           `json:"-" openapi:"ignored"`
	IsMainAddressAsBool     bool                           `json:"-" openapi:"ignored"`
}

func (data *AddressesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"staffId":           govalidity.New("staffId"),
		"customerId":        govalidity.New("customerId"),
		"city":              govalidity.New("city").Optional(),
		"street":            govalidity.New("street").MinLength(2),
		"name":              govalidity.New("name").MinMaxLength(2, 25).Required(),
		"postalCode":        govalidity.New("postalCode").MinMaxLength(2, 25),
		"buildingNumber":    govalidity.New("buildingNumber").MinMaxLength(2, 25).Required(),
		"isDeliveryAddress": govalidity.New("isDeliveryAddress").In([]string{"true", "false"}).Required(),
		"isMainAddress":     govalidity.New("isMainAddress").In([]string{"true", "false"}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// if staffId and customerId are both nil, return error
	if data.StaffID == nil && data.CustomerID == nil {
		return govalidity.ValidityResponseErrors{
			"staffId":    []string{"staffId or customerId is required"},
			"customerId": []string{"staffId or customerId is required"},
		}
	}

	// if staffId and customerId are both not nil, return error
	if data.StaffID != nil && data.CustomerID != nil {
		return govalidity.ValidityResponseErrors{
			"staffId":    []string{"only one of staffId or customerId is allowed"},
			"customerId": []string{"only one of staffId or customerId is allowed"},
		}
	}

	// Check isDeliveryAddress
	if data.IsDeliveryAddress == "true" {
		data.IsDeliveryAddressAsBool = true
	} else {
		data.IsDeliveryAddressAsBool = false
	}

	// Check isMainAddress
	if data.IsMainAddress == "true" {
		data.IsMainAddressAsBool = true
	} else {
		data.IsMainAddressAsBool = false
	}

	return nil
}
