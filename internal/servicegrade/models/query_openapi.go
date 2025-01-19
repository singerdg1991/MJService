package models

import "github.com/hoitek/Maja-Service/internal/servicegrade/domain"

/*
 * @apiDefine: ServiceGradesQueryResponseDataItem
 */
type ServiceGradesQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
	Grade       int     `json:"grade" openapi:"example:0;required"`
	Color       string  `json:"color" openapi:"example:#000000;required"`
}

/*
 * @apiDefine: ServiceGradesQueryResponseData
 */
type ServiceGradesQueryResponseData struct {
	Limit      int                                  `json:"limit" openapi:"example:10"`
	Offset     int                                  `json:"offset" openapi:"example:0"`
	Page       int                                  `json:"page" openapi:"example:1"`
	TotalRows  int                                  `json:"totalRows" openapi:"example:1"`
	TotalPages int                                  `json:"totalPages" openapi:"example:1"`
	Items      []ServiceGradesQueryResponseDataItem `json:"items" openapi:"$ref:ServiceGradesQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: ServiceGradesQueryResponse
 */
type ServiceGradesQueryResponse struct {
	StatusCode int                            `json:"statusCode" openapi:"example:200;"`
	Data       ServiceGradesQueryResponseData `json:"data" openapi:"$ref:ServiceGradesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: ServiceGradesQueryNotFoundResponse
 */
type ServiceGradesQueryNotFoundResponse struct {
	ServiceGrades []domain.ServiceGrade `json:"servicegrades" openapi:"$ref:ServiceGrade;type:array"`
}
