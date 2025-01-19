package models

import (
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: GracesCreateRequestBody
 */
type GracesCreateRequestBody struct {
	IsAutoRewardSetEnable       string              `json:"isAutoRewardSetEnable" openapi:"example:true"`
	GraceNumber                 int                 `json:"graceNumber" openapi:"example:1"`
	Title                       string              `json:"title" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description                 string              `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	RewardID                    uint                `json:"rewardId" openapi:"example:1"`
	Reward                      *domain.GraceReward `json:"-" openapi:"ignored"`
	IsAutoRewardSetEnableAsBool bool                `json:"-" openapi:"ignored"`
}

func (data *GracesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"isAutoRewardSetEnable": govalidity.New("isAutoRewardSetEnable").In([]string{"true", "false"}).Required(),
		"graceNumber":           govalidity.New("graceNumber").Required(),
		"title":                 govalidity.New("title").MinMaxLength(3, 25).Required(),
		"description":           govalidity.New("description"),
		"rewardId":              govalidity.New("rewardId").Int().Min(1).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Convert string to bool
	if data.IsAutoRewardSetEnable == "true" {
		data.IsAutoRewardSetEnableAsBool = true
	} else {
		data.IsAutoRewardSetEnableAsBool = false
	}

	return nil
}
