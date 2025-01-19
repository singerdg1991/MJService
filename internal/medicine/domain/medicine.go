package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: Medicine
 */
type Medicine struct {
	ID           uint       `json:"id" openapi:"example:1"`
	Name         string     `json:"name" openapi:"example:name;required"`
	Code         *string    `json:"code" openapi:"example:code;required"`
	Availability *string    `json:"availability" openapi:"example:availability;required"`
	Manufacturer *string    `json:"manufacturer" openapi:"example:manufacturer;required"`
	PurposeOfUse *string    `json:"purposeOfUse" openapi:"example:purposeOfUse;required"`
	Instruction  *string    `json:"instruction" openapi:"example:instruction;required"`
	SideEffects  *string    `json:"sideEffects" openapi:"example:sideEffects;required"`
	Conditions   *string    `json:"conditions" openapi:"example:conditions;required"`
	Description  *string    `json:"description" openapi:"example:description;required"`
	CreatedAt    time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt    time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt    *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func (u *Medicine) TableName() string {
	return "medicines"
}

func (u *Medicine) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
