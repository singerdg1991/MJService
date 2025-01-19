package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Maja-Service/internal/ability/domain"
)

/*
 * @apiDefine: AbilitiesCreateRequestBody
 */
type AbilitiesCreateRequestBody struct {
	Name string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *AbilitiesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
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
 * @apiDefine: AbilitiesCreateResponse
 */
type AbilitiesCreateResponse struct {
	Ability domain.Ability `json:"ability" openapi:"$ref:Ability;type:object;"`
}
