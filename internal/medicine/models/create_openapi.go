package models

/*
 * @apiDefine: MedicinesResponseData
 */
type MedicinesResponseData struct {
	ID           uint   `json:"id" openapi:"example:1"`
	Name         string `json:"name" openapi:"example:saeed"`
	Code         string `json:"code" openapi:"example:code"`
	Availability string `json:"availability" openapi:"example:availability"`
	Manufacturer string `json:"manufacturer" openapi:"example:manufacturer"`
	PurposeOfUse string `json:"purposeOfUse" openapi:"example:purposeOfUse"`
	Instruction  string `json:"instruction" openapi:"example:instruction"`
	SideEffects  string `json:"sideEffects" openapi:"example:sideEffects"`
	Conditions   string `json:"conditions" openapi:"example:conditions"`
	Description  string `json:"description" openapi:"example:description"`
}

/*
 * @apiDefine: MedicinesCreateResponse
 */
type MedicinesCreateResponse struct {
	StatusCode int                   `json:"statusCode" openapi:"example:200"`
	Data       MedicinesResponseData `json:"data" openapi:"$ref:MedicinesResponseData"`
}
