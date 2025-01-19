package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Maja-Service/internal/city/domain"
)

/*
 * @apiDefine: CitiesCreateRequestBody
 */
type CitiesCreateRequestBody struct {
	Name string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *CitiesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name": govalidity.New("name").MinMaxLength(3, 25).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CitiesCreateResponse
 */
type CitiesCreateResponse struct {
	City domain.City `json:"city" openapi:"$ref:City;type:object;"`
}
