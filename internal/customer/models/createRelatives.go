package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersCreateRelativesRequestBody
 */
type CustomersCreateRelativesRequestBody struct {
	CustomerID            int    `json:"customerId" openapi:"example:1"`
	CityID                int    `json:"cityId" openapi:"example:1"`
	FirstName             string `json:"firstName" openapi:"example:John;required"`
	LastName              string `json:"lastName" openapi:"example:Doe;required"`
	PhoneNumber           string `json:"phoneNumber" openapi:"example:+989204005707"`
	Relation              string `json:"relation" openapi:"example:father"`
	AddressName           string `json:"addressName" openapi:"example:home"`
	AddressStreet         string `json:"addressStreet" openapi:"example:street"`
	AddressBuildingNumber string `json:"addressBuildingNumber" openapi:"example:1"`
	AddressPostalCode     string `json:"addressPostalCode" openapi:"example:1234567890"`
}

func (data *CustomersCreateRelativesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"customerId":            govalidity.New("customerId").Int().Min(1).Required(),
		"cityId":                govalidity.New("cityId").Int().Min(1).Required(),
		"firstName":             govalidity.New("firstName").Required(),
		"lastName":              govalidity.New("lastName").Optional(),
		"phoneNumber":           govalidity.New("phoneNumber").Required(),
		"relation":              govalidity.New("relation").Required(),
		"addressName":           govalidity.New("addressName").Required(),
		"addressStreet":         govalidity.New("addressStreet").Required(),
		"addressBuildingNumber": govalidity.New("addressBuildingNumber").Required(),
		"addressPostalCode":     govalidity.New("addressPostalCode").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
