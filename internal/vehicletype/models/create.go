package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Maja-Service/internal/vehicletype/domain"
)

/*
 * @apiDefine: VehicleTypesCreateRequestBody
 */
type VehicleTypesCreateRequestBody struct {
	Name string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *VehicleTypesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
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
 * @apiDefine: VehicleTypesCreateResponse
 */
type VehicleTypesCreateResponse struct {
	VehicleType domain.VehicleType `json:"vehicletype" openapi:"$ref:VehicleType;type:object;"`
}
