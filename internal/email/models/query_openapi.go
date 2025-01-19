package models

import "github.com/hoitek/Maja-Service/internal/email/domain"

/*
 * @apiDefine: EmailsQueryResponseData
 */
type EmailsQueryResponseData struct {
	Limit      int                  `json:"limit" openapi:"example:10"`
	Offset     int                  `json:"offset" openapi:"example:0"`
	Page       int                  `json:"page" openapi:"example:1"`
	TotalRows  int                  `json:"totalRows" openapi:"example:1"`
	TotalPages int                  `json:"totalPages" openapi:"example:1"`
	Items      []EmailsResponseData `json:"items" openapi:"$ref:EmailsResponseData;type:array"`
}

/*
 * @apiDefine: EmailsQueryResponse
 */
type EmailsQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200;"`
	Data       EmailsQueryResponseData `json:"data" openapi:"$ref:EmailsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: EmailsQueryNotFoundResponse
 */
type EmailsQueryNotFoundResponse struct {
	Emails []domain.Email `json:"emails" openapi:"$ref:Email;type:array"`
}
