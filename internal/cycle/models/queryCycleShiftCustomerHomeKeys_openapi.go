package models

import "github.com/hoitek/Maja-Service/internal/cycle/domain"

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysResponseData
 */
type CyclesQueryShiftCustomerHomeKeysResponseData struct {
	Limit      int                                            `json:"limit" openapi:"example:10"`
	Offset     int                                            `json:"offset" openapi:"example:0"`
	Page       int                                            `json:"page" openapi:"example:1"`
	TotalRows  int                                            `json:"totalRows" openapi:"example:1"`
	TotalPages int                                            `json:"totalPages" openapi:"example:1"`
	Items      []CyclesCreateShiftCustomerHomeKeyResponseData `json:"items" openapi:"$ref:CyclesCreateShiftCustomerHomeKeyResponseData;type:array"`
}

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysResponse
 */
type CyclesQueryShiftCustomerHomeKeysResponse struct {
	StatusCode int                                          `json:"statusCode" openapi:"example:200"`
	Data       CyclesQueryShiftCustomerHomeKeysResponseData `json:"data" openapi:"$ref:CyclesQueryShiftCustomerHomeKeysResponseData"`
}

/*
 * @apiDefine: CyclesQueryShiftCustomerHomeKeysNotFoundResponse
 */
type CyclesQueryShiftCustomerHomeKeysNotFoundResponse struct {
	Cycles []domain.Cycle `json:"cycles" openapi:"$ref:Cycle;type:array"`
}
