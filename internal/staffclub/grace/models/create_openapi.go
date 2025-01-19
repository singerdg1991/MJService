package models

import "github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"

/*
 * @apiDefine: GracesResponseData
 */
type GracesResponseData struct {
	ID                    uint               `json:"id" openapi:"example:1"`
	RewardID              uint               `json:"rewardId" openapi:"example:1"`
	Reward                domain.GraceReward `json:"reward" openapi:"$ref:GraceReward"`
	IsAutoRewardSetEnable bool               `json:"isAutoRewardSetEnable" openapi:"example:true"`
	GraceNumber           int64              `json:"graceNumber" openapi:"example:1"`
	Title                 string             `json:"title" openapi:"example:test"`
	Description           string             `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: GracesCreateResponse
 */
type GracesCreateResponse struct {
	StatusCode int                `json:"statusCode" openapi:"example:200"`
	Data       GracesResponseData `json:"data" openapi:"$ref:GracesResponseData"`
}
