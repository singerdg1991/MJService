package domain

import (
	"encoding/json"
	"time"
)

type StaffStaffTypes struct {
	ID          int64      `json:"id"`
	StaffID     int64      `json:"staffId"`
	ShiftTypeID int64      `json:"shiftTypeId"`
	CreatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffStaffTypes() *StaffStaffTypes {
	return &StaffStaffTypes{}
}

func (ns *StaffStaffTypes) TableName() string {
	return "staffStaffTypes"
}

func (ns *StaffStaffTypes) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
