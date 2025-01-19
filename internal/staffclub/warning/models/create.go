package models

import (
	"net/http"

	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: WarningsCreateRequestBody
 */
type WarningsCreateRequestBody struct {
	IsAutoRewardSetEnable       string                    `json:"isAutoRewardSetEnable" openapi:"example:true"`
	WarningNumber               int                       `json:"warningNumber" openapi:"example:1"`
	Title                       string                    `json:"title" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Description                 string                    `json:"description" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	PunishmentID                uint                      `json:"punishmentId" openapi:"example:1"`
	Punishment                  *domain.WarningPunishment `json:"-" openapi:"ignored"`
	IsAutoRewardSetEnableAsBool bool                      `json:"-" openapi:"ignored"`
}

func (data *WarningsCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"isAutoRewardSetEnable": govalidity.New("isAutoRewardSetEnable").In([]string{"true", "false"}).Required(),
		"warningNumber":         govalidity.New("warningNumber").Required(),
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
