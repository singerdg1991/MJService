package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: UsersResetPasswordRequestRequestBody
 */
type UsersResetPasswordRequestRequestBody struct {
	UserID int `json:"userId" openapi:"example:1;required;"`
}

func (data *UsersResetPasswordRequestRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userId": govalidity.New("email").Int().Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: UsersResetPasswordRequestResponseData
 */
type UsersResetPasswordRequestResponseData struct {
	Message string `json:"message" openapi:"example:Password reset successfully;"`
}

/*
 * @apiDefine: UsersResetPasswordRequestResponse
 */
type UsersResetPasswordRequestResponse struct {
	StatusCode int                                   `json:"statusCode" openapi:"example:200;"`
	Data       UsersResetPasswordRequestResponseData `json:"data" openapi:"$ref:UsersResetPasswordRequestResponseData;type:object;"`
}
