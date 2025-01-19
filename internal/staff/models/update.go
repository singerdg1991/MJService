package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: StaffsUpdateContractRequestParams
 */
type StaffsUpdateRequestParams struct {
	ID int `json:"id,string" openapi:"example:id;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *StaffsUpdateRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id": govalidity.New("id").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: StaffsUpdateRequestBody
 */
type StaffsUpdateRequestBody struct {
	ID int `json:"id" openapi:"ignored"`
}

func (data *StaffsUpdateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"joinedAt":                govalidity.New("joinedAt").Optional(),
		"contractStartedAt":       govalidity.New("contractStartedAt").Optional(),
		"contractExpiresAt":       govalidity.New("contractExpiresAt").Optional(),
		"availableShifts":         govalidity.New("availableShifts").Optional(),
		"trialTime":               govalidity.New("trialTime").Optional(),
		"listOfContract":          govalidity.New("listOfContract").Optional(),
		"organizationNumber":      govalidity.New("organizationNumber").Optional(),
		"jobTitle":                govalidity.New("jobTitle").Optional(),
		"percentLengthInContract": govalidity.New("percentLengthInContract").Optional(),
		"hourLengthInContract":    govalidity.New("hourLengthInContract").Optional(),
		"salary":                  govalidity.New("salary").Optional(),
		"sections":                govalidity.New("sections").Optional(),
		"teams":                   govalidity.New("teams").Optional(),
		"abilities":               govalidity.New("abilities").Optional(),
		"limitations":             govalidity.New("limitations").Optional(),
		"role":                    govalidity.New("role").Optional(),
		"staffType":               govalidity.New("staffType").Optional(),
		"experienceAmount":        govalidity.New("experienceAmount").Optional(),
		"experienceAmountUnit":    govalidity.New("experienceAmountUnit").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
