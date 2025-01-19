package models

import (
	"net/http"

	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: AttentionsCreateRequestBody
 */
type AttentionsCreateRequestBody struct {
	IsAutoRewardSetEnable       string                      `json:"isAutoRewardSetEnable" openapi:"example:true"`
	AttentionNumber             int                         `json:"attentionNumber" openapi:"example:1"`
	Title                       string                      `json:"title" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description                 string                      `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	PunishmentID                uint                        `json:"punishmentId" openapi:"example:1"`
	Punishment                  *domain.AttentionPunishment `json:"-" openapi:"ignored"`
	IsAutoRewardSetEnableAsBool bool                        `json:"-" openapi:"ignored"`
}

func (data *AttentionsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"isAutoRewardSetEnable": govalidity.New("isAutoRewardSetEnable").In([]string{"true", "false"}).Required(),
		"attentionNumber":       govalidity.New("attentionNumber").Required(),
		"title":                 govalidity.New("title").MinMaxLength(3, 25).Required(),
		"description":           govalidity.New("description"),
		"punishmentId":          govalidity.New("punishmentId").Int().Min(1).Required(),
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
