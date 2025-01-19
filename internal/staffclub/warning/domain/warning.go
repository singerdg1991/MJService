package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: WarningPunishment
 */
type WarningPunishment struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: Warning
 */
type Warning struct {
	ID                    uint               `json:"id" openapi:"example:1"`
	PunishmentID          uint               `json:"punishmentId" openapi:"example:1"`
	Punishment            *WarningPunishment `json:"punishment" openapi:"$ref:WarningPunishment;type:object;"`
	IsAutoRewardSetEnable bool               `json:"isAutoRewardSetEnable" openapi:"example:true"`
	WarningNumber         int                `json:"warningNumber" openapi:"example:1"`
	Title                 string             `json:"title" openapi:"example:John;required"`
	Description           *string            `json:"description" openapi:"example:John;required"`
	CreatedAt             time.Time          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             time.Time          `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             *time.Time         `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Warning) TableName() string {
	return "staffClubWarnings"
}

func (u *Warning) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
