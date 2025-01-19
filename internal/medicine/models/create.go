package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"net/http"
)

/*
 * @apiDefine: MedicinesCreateRequestBody
 */
type MedicinesCreateRequestBody struct {
	Name         string `json:"name" openapi:"example:name;required;maxLen:100;minLen:2;"`
	Code         string `json:"code" openapi:"example:code;required;maxLen:100;minLen:2;"`
	Availability string `json:"availability" openapi:"example:availability;required;maxLen:100;minLen:2;"`
	Manufacturer string `json:"manufacturer" openapi:"example:manufacturer;required;maxLen:100;minLen:2;"`
	PurposeOfUse string `json:"purposeOfUse" openapi:"example:purposeOfUse;required;maxLen:100;minLen:2;"`
	Instruction  string `json:"instruction" openapi:"example:instruction;required;maxLen:100;minLen:2;"`
	SideEffects  string `json:"sideEffects" openapi:"example:sideEffects;required;maxLen:100;minLen:2;"`
	Conditions   string `json:"conditions" openapi:"example:conditions;required;maxLen:100;minLen:2;"`
	Description  string `json:"description" openapi:"example:description;required;maxLen:100;minLen:2;"`
}

func (data *MedicinesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":         govalidity.New("name").MinMaxLength(3, 255).Required(),
		"code":         govalidity.New("code").MinMaxLength(3, 255).Required(),
		"availability": govalidity.New("availability").MinMaxLength(3, 2000).Optional(),
		"manufacturer": govalidity.New("manufacturer").MinMaxLength(3, 2000).Optional(),
		"purposeOfUse": govalidity.New("purposeOfUse").MinMaxLength(3, 2000).Optional(),
		"instruction":  govalidity.New("instruction").MinMaxLength(3, 2000).Optional(),
		"sideEffects":  govalidity.New("sideEffects").MinMaxLength(3, 2000).Optional(),
		"conditions":   govalidity.New("conditions").MinMaxLength(3, 2000).Optional(),
		"description":  govalidity.New("description").MinMaxLength(3, 2000).Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
