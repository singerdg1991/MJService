package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"net/http"
)

/*
 * @apiDefine: UsersResetPasswordRequestBody
 */
type UsersResetPasswordRequestBody struct {
	Email        string `json:"email" openapi:"example:sgh370@yahoo.com;required;"`
	NewPassword  string `json:"newPassword" openapi:"example:111111;required;"`
	ExchangeCode string `json:"exchangeCode" openapi:"example:123456;required;maxLen:100;minLen:2;"`
}

func (data *UsersResetPasswordRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"email":        govalidity.New("email").Email().Required(),
		"newPassword":  govalidity.New("newPassword").Required(),
		"exchangeCode": govalidity.New("exchangeCode").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: UsersResetPassword
 */
type UsersResetPassword struct {
	Status string `json:"status" openapi:"example:success;"`
}

/*
 * @apiDefine: UsersResetPasswordResponse
 */
type UsersResetPasswordResponse struct {
	User domain.User `json:"user" openapi:"$ref:User;type:object;"`
}
