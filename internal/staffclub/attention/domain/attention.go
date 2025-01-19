package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: AttentionPunishment
 */
type AttentionPunishment struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: Attention
 */
type Attention struct {
	ID                    uint                 `json:"id" openapi:"example:1"`
	PunishmentID          uint                 `json:"punishmentId" openapi:"example:1"`
	Punishment            *AttentionPunishment `json:"punishment" openapi:"$ref:AttentionPunishment;type:object;"`
	IsAutoRewardSetEnable bool                 `json:"isAutoRewardSetEnable" openapi:"example:true"`
	AttentionNumber       int                  `json:"attentionNumber" openapi:"example:1"`
	Title                 string               `json:"title" openapi:"example:John;required"`
	Description           *string              `json:"description" openapi:"example:John;required"`
	CreatedAt             time.Time            `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             time.Time            `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             *time.Time           `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Attention) TableName() string {
	return "staffClubAttentions"
}

func (u *Attention) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
