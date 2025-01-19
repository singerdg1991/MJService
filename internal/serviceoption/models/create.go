package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ServiceOptionsCreateRequestBody
 */
type ServiceOptionsCreateRequestBody struct {
	ServiceTypeID int64  `json:"serviceTypeId" openapi:"example:1;required;"`
	Name          string `json:"name" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description   string `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
}

func (data *ServiceOptionsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"serviceTypeId": govalidity.New("serviceTypeId").Int().Required(),
		"name":          govalidity.New("name").MinMaxLength(3, 25).Required(),
		"description":   govalidity.New("description").Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
