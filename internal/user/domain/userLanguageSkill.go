package domain

import (
	"encoding/json"
	"time"
)

type UserLanguageSkill struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"userId"`
	LanguageSkillID int64      `json:"languageSkillId"`
	CreatedAt       time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt       *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewUserLanguageSkill() *UserLanguageSkill {
	return &UserLanguageSkill{}
}

func (ns *UserLanguageSkill) TableName() string {
	return "userLanguageSkills"
}

func (ns *UserLanguageSkill) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
