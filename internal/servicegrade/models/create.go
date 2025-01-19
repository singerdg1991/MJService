package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ServiceGradesCreateRequestBody
 */
type ServiceGradesCreateRequestBody struct {
	Name        string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description string `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Grade       int    `json:"grade" openapi:"example:0;required"`
	Color       string `json:"color" openapi:"example:#000000;required"`
}

func (data *ServiceGradesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":        govalidity.New("name").MinMaxLength(3, 25).Required(),
		"description": govalidity.New("description").Optional(),
		"grade":       govalidity.New("grade").Int().Min(0).Required(),
		"color":       govalidity.New("color").HexColor().Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
