package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: CustomersCreateCreditDetailsRequestBody
 */
type CustomersCreateCreditDetailsRequestBody struct {
	CustomerID        int64  `json:"-" openapi:"ignored"`
	UserID            int64  `json:"userId" openapi:"example:1"`
	BillingAddressID  int64  `json:"billingAddressId" openapi:"example:1"`
	BankAccountNumber string `json:"bankAccountNumber" openapi:"example:123456789"`
}

func (data *CustomersCreateCreditDetailsRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userId":            govalidity.New("userId").Int().Min(1).Required(),
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
