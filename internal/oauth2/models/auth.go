package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"net/http"
)

/*
 * @apiDefine: OAuth2AuthRequestBody
 */
type OAuth2AuthRequestBody struct {
	Username     string `json:"username" openapi:"example:saeed@gmail.com;required;maxLen:100;minLen:2;"`
	Password     string `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Type         string `json:"type" openapi:"example:email;required;maxLen:100;minLen:2;"`
	ExchangeCode string `json:"exchangeCode" openapi:"example:123456;required;maxLen:100;minLen:2;"`
}

func (data *OAuth2AuthRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"username":     govalidity.New("username").MinMaxLength(3, 100).Required(),
		"password":     govalidity.New("password").MinMaxLength(3, 100).Required(),
		"type":         govalidity.New("type").Required(),
		"exchangeCode": govalidity.New("exchangeCode"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: OAuth2Response
 */
type OAuth2Response struct {
	AccessToken  string      `json:"accessToken" openapi:"example:xxxxxxxxxxxxxxxxxxxxxxxx;"`
	RefreshToken string      `json:"refreshToken" openapi:"example:xxxxxxxxxxxxxxxxxxxxxxxx;"`
	User         domain.User `json:"user" openapi:"$ref:User;type:object;"`
}

/*
 * @apiDefine: OAuth2AuthResponse
 */
type OAuth2AuthResponse struct {
	OAuth2 OAuth2Response `json:"oauth2" openapi:"$ref:OAuth2Response;type:object;"`
}
