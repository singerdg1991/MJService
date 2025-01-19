package models

import "github.com/hoitek/Maja-Service/internal/prescription/domain"

/*
 * @apiDefine: PrescriptionsQueryResponseData
 */
type PrescriptionsQueryResponseData struct {
	Limit      int                               `json:"limit" openapi:"example:10"`
	Offset     int                               `json:"offset" openapi:"example:0"`
	Page       int                               `json:"page" openapi:"example:1"`
	TotalRows  int                               `json:"totalRows" openapi:"example:1"`
	TotalPages int                               `json:"totalPages" openapi:"example:1"`
	Items      []PrescriptionsCreateResponseData `json:"items" openapi:"$ref:PrescriptionsCreateResponseData;type:array"`
}

/*
 * @apiDefine: PrescriptionsQueryResponse
 */
type PrescriptionsQueryResponse struct {
	StatusCode int                            `json:"statusCode" openapi:"example:200"`
	Data       PrescriptionsQueryResponseData `json:"data" openapi:"$ref:PrescriptionsQueryResponseData;type:object"`
}

/*
 * @apiDefine: PrescriptionsQueryNotFoundResponse
 */
type PrescriptionsQueryNotFoundResponse struct {
	Prescriptions []domain.Prescription `json:"prescriptions" openapi:"$ref:Prescription;type:array"`
}
