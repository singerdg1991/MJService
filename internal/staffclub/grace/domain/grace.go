package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: GraceReward
 */
type GraceReward struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:John;required"`
	Description string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: Grace
 */
type Grace struct {
	ID                    uint         `json:"id" openapi:"example:1"`
	RewardID              uint         `json:"rewardId" openapi:"example:1"`
	Reward                *GraceReward `json:"reward" openapi:"$ref:GraceReward;type:object;"`
	IsAutoRewardSetEnable bool         `json:"isAutoRewardSetEnable" openapi:"example:true"`
	GraceNumber           int          `json:"graceNumber" openapi:"example:1"`
	Title                 string       `json:"title" openapi:"example:John;required"`
	Description           *string      `json:"description" openapi:"example:John;required"`
	CreatedAt             time.Time    `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt             time.Time    `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt             *time.Time   `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Grace) TableName() string {
	return "staffClubGraces"
}

func (u *Grace) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
