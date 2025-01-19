package models

import "github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"

/*
 * @apiDefine: HolidaysQueryResponseData
 */
type HolidaysQueryResponseData struct {
	Limit      int                    `json:"limit" openapi:"example:10"`
	Offset     int                    `json:"offset" openapi:"example:0"`
	Page       int                    `json:"page" openapi:"example:1"`
	TotalRows  int                    `json:"totalRows" openapi:"example:1"`
	TotalPages int                    `json:"totalPages" openapi:"example:1"`
	Items      []HolidaysResponseData `json:"items" openapi:"$ref:HolidaysResponseData;type:array"`
}

/*
 * @apiDefine: HolidaysQueryResponse
 */
type HolidaysQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200;"`
	Data       HolidaysQueryResponseData `json:"data" openapi:"$ref:HolidaysQueryResponseData;type:object;"`
}

/*
 * @apiDefine: HolidaysQueryNotFoundResponse
 */
type HolidaysQueryNotFoundResponse struct {
	Holidays []domain.Holiday `json:"holidays" openapi:"$ref:Holiday;type:array"`
}
