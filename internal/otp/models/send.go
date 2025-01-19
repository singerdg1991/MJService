package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/otp"
	"net/http"
)

/*
 * @apiDefine: OTPSendRequest
 */
type OTPSendRequest struct {
	Username string `json:"username" openapi:"example:sgh370@yahoo.com;required;"`
	Password string `json:"password" openapi:"example:111111;required;"`
	Type     string `json:"type" openapi:"example:email;required;"`
	Reason   string `json:"reason" openapi:"example:login;required;"`
}

func (data *OTPSendRequest) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"username": govalidity.New("username").Required(),
		"password": govalidity.New("password").Optional(),
		"type":     govalidity.New("type").In([]string{otp.TypeEmail, otp.TypePhone}).Required(),
		"reason":   govalidity.New("reason").In([]string{otp.ReasonRegister, otp.ReasonChangePassword, otp.ReasonLogin, otp.ReasonResetPassword}).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: OTPSendResponse
 */
type OTPSendResponse struct {
	Token    string `json:"token"`
	CoolDown int    `json:"coolDown"`
	Code     string `json:"code"`
}
