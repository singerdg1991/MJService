package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsQueryProfileRequestParams
 */
type StaffsQueryProfileRequestParams struct {
	ID      int `json:"id,string,omitempty" openapi:"example:1"`
	StaffID int `json:"staffId,string,omitempty" openapi:"example:1;required"`
}

func (data *StaffsQueryProfileRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":      govalidity.New("id").Int().Optional(),
		"staffId": govalidity.New("staffId").Int().Required(),
	}

	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
