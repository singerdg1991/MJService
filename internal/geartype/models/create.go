package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Maja-Service/internal/geartype/domain"
)

/*
 * @apiDefine: GearTypesCreateRequestBody
 */
type GearTypesCreateRequestBody struct {
	Name string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *GearTypesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
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
 * @apiDefine: GearTypesCreateResponse
 */
type GearTypesCreateResponse struct {
	GearType domain.GearType `json:"geartype" openapi:"$ref:GearType;type:object;"`
}
