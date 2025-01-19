package models

import "github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"

/*
 * @apiDefine: AttentionsResponseData
 */
type AttentionsResponseData struct {
	ID                    uint                       `json:"id" openapi:"example:1"`
	PunishmentID          uint                       `json:"punishmentId" openapi:"example:1"`
	Punishment            domain.AttentionPunishment `json:"punishment" openapi:"$ref:AttentionPunishment"`
	IsAutoRewardSetEnable bool                       `json:"isAutoRewardSetEnable" openapi:"example:true"`
	AttentionNumber       int64                      `json:"attentionNumber" openapi:"example:1"`
	Title                 string                     `json:"title" openapi:"example:test"`
	Description           string                     `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: AttentionsCreateResponse
 */
type AttentionsCreateResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200"`
	Data       AttentionsResponseData `json:"data" openapi:"$ref:AttentionsResponseData"`
}
