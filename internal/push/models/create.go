package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: PushesCreateRequestBody
 */
type PushesCreateRequestBody struct {
	UserID     int    `json:"userId" openapi:"example:1;required;"`
	Endpoint   string `json:"endpoint" openapi:"example:endpoint;required;"`
	KeysAuth   string `json:"keysAuth" openapi:"example:keysAuth;required;"`
	KeysP256dh string `json:"keysP256dh" openapi:"example:keysP256dh;required;"`
}

func (data *PushesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userId":     govalidity.New("userId").Int().Min(1).Required(),
		"endpoint":   govalidity.New("endpoint").MaxLength(1000).Required(),
		"keysAuth":   govalidity.New("keysAuth").MaxLength(1000).Required(),
		"keysP256dh": govalidity.New("keysP256dh").MaxLength(1000).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
