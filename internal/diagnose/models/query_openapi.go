package models

import "github.com/hoitek/Maja-Service/internal/diagnose/domain"

/*
 * @apiDefine: DiagnosesQueryResponseData
 */
type DiagnosesQueryResponseData struct {
	Limit      int                           `json:"limit" openapi:"example:10"`
	Offset     int                           `json:"offset" openapi:"example:0"`
	Page       int                           `json:"page" openapi:"example:1"`
	TotalRows  int                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                           `json:"totalPages" openapi:"example:1"`
	Items      []DiagnosesCreateResponseData `json:"items" openapi:"$ref:DiagnosesCreateResponseData;type:array"`
}

/*
 * @apiDefine: DiagnosesQueryResponse
 */
type DiagnosesQueryResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       DiagnosesQueryResponseData `json:"data" openapi:"$ref:DiagnosesQueryResponseData;type:object"`
}

/*
 * @apiDefine: DiagnosesQueryNotFoundResponse
 */
type DiagnosesQueryNotFoundResponse struct {
	Diagnoses []domain.Diagnose `json:"diagnoses" openapi:"$ref:Diagnose;type:array"`
}
