package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: OTPVerifyRequest
 */
type OTPVerifyRequest struct {
	Token string `json:"token" openapi:"example:xxxxxxxxxxxxxxxxxxxxxxxxxxxx;required;"`
	Code  string `json:"code" openapi:"example:111111;required;"`
}

func (data *OTPVerifyRequest) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"token": govalidity.New("token").Required(),
		"code":  govalidity.New("code").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: OTPVerifyResponse
 */
type OTPVerifyResponse struct {
	ExchangeCode string `json:"exchangeCode"`
}
