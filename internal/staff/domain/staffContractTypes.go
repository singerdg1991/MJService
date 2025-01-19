package domain

import (
	"encoding/json"
	"time"
)

type StaffContractTypes struct {
	ID             int64      `json:"id"`
	StaffID        int64      `json:"staffId"`
	ContractTypeID int64      `json:"contractTypeId"`
	CreatedAt      time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt      *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffContractTypes() *StaffContractTypes {
	return &StaffContractTypes{}
}

func (ns *StaffContractTypes) TableName() string {
	return "staffContractTypes"
}

func (ns *StaffContractTypes) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
