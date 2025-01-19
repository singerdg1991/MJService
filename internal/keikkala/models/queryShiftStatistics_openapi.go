package models

import "github.com/hoitek/Maja-Service/internal/keikkala/domain"

/*
 * @apiDefine: KeikkalasQueryShiftStatisticsResponseData
 */
type KeikkalasQueryShiftStatisticsResponseData struct {
	MorningCount int `json:"morningCount" openapi:"example:1;"`
	EveningCount int `json:"eveningCount" openapi:"example:1;"`
	NightCount   int `json:"nightCount" openapi:"example:1;"`
}

/*
 * @apiDefine: KeikkalasQueryShiftStatisticsResponse
 */
type KeikkalasQueryShiftStatisticsResponse struct {
	StatusCode int                                       `json:"statusCode" openapi:"example:200;"`
	Data       KeikkalasQueryShiftStatisticsResponseData `json:"data" openapi:"$ref:KeikkalasQueryShiftStatisticsResponseData;type:object;"`
}

/*
 * @apiDefine: KeikkalasQueryShiftStatisticsNotFoundResponse
 */
type KeikkalasQueryShiftStatisticsNotFoundResponse struct {
	Keikkalas []domain.Keikkala `json:"keikkalas" openapi:"$ref:Keikkala;type:array"`
}
