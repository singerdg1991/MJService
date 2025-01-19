package models

import "github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"

/*
 * @apiDefine: WarningsResponseData
 */
type WarningsResponseData struct {
	ID                    uint                     `json:"id" openapi:"example:1"`
	PunishmentID          uint                     `json:"punishmentId" openapi:"example:1"`
	Punishment            domain.WarningPunishment `json:"punishment" openapi:"$ref:WarningPunishment"`
	IsAutoRewardSetEnable bool                     `json:"isAutoRewardSetEnable" openapi:"example:true"`
	WarningNumber         int64                    `json:"warningNumber" openapi:"example:1"`
	Title                 string                   `json:"title" openapi:"example:test"`
	Description           string                   `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: WarningsCreateResponse
 */
type WarningsCreateResponse struct {
	StatusCode int                  `json:"statusCode" openapi:"example:200"`
	Data       WarningsResponseData `json:"data" openapi:"$ref:WarningsResponseData"`
}
