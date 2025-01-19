package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: OAuth2RefreshRequestBody
 */
type OAuth2RefreshRequestBody struct {
	RefreshToken string `json:"refreshToken" openapi:"example:xxxxxxxxxxxxxxxxxxxxxxxxx;required;"`
}

func (data *OAuth2RefreshRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"refreshToken": govalidity.New("username").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: OAuth2RefreshResponse
 */
type OAuth2RefreshResponse struct {
	OAuth2 OAuth2Response `json:"oauth2" openapi:"$ref:OAuth2Response;type:object;"`
}
