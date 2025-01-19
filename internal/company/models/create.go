package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"

	"github.com/hoitek/Maja-Service/internal/company/domain"
)

/*
 * @apiDefine: CompaniesCreateRequestBody
 */
type CompaniesCreateRequestBody struct {
	Name string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *CompaniesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
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
 * @apiDefine: CompaniesCreateResponse
 */
type CompaniesCreateResponse struct {
	Company domain.Company `json:"company" openapi:"$ref:Company;type:object;"`
}
