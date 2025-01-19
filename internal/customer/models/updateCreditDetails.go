package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersUpdateCreditDetailsRequestParams
 */
type CustomersUpdateCreditDetailsRequestParams struct {
	ID         int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
	CustomerID int `json:"customerid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CustomersUpdateCreditDetailsRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":         govalidity.New("id").Int().Required(),
		"customerid": govalidity.New("customerId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CustomersUpdateCreditDetailsRequestBody
 */
type CustomersUpdateCreditDetailsRequestBody struct {
	ID                int64  `json:"-" openapi:"ignored"`
	CustomerID        int64  `json:"-" openapi:"ignored"`
	BillingAddressID  int64  `json:"billingAddressId" openapi:"example:1"`
	BankAccountNumber string `json:"bankAccountNumber" openapi:"example:123456789"`
}

func (data *CustomersUpdateCreditDetailsRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"billingAddressId":  govalidity.New("billingAddressId").Int().Min(1).Required(),
		"bankAccountNumber": govalidity.New("bankAccountNumber").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
