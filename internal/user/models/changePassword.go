package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"net/http"
)

/*
 * @apiDefine: UsersChangePasswordRequestBody
 */
type UsersChangePasswordRequestBody struct {
	CurrentPassword string `json:"currentPassword" openapi:"example:111111;required;"`
	NewPassword     string `json:"newPassword" openapi:"example:111111;required;"`
}

func (data *UsersChangePasswordRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"currentPassword": govalidity.New("currentPassword").Required(),
		"newPassword":     govalidity.New("newPassword").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: UsersChangePassword
 */
type UsersChangePassword struct {
	Status string `json:"status" openapi:"example:success;"`
}

/*
 * @apiDefine: UsersChangePasswordResponse
 */
type UsersChangePasswordResponse struct {
	User domain.User `json:"user" openapi:"$ref:User;type:object;"`
}
