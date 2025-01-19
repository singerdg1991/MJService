package domain

import (
	"encoding/json"
	"time"
)

type StaffShiftTypes struct {
	ID          int64      `json:"id"`
	StaffID     int64      `json:"staffId"`
	ShiftTypeID int64      `json:"shiftTypeId"`
	CreatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"-" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffShiftTypes() *StaffShiftTypes {
	return &StaffShiftTypes{}
}

func (ns *StaffShiftTypes) TableName() string {
	return "staffShiftTypes"
}

func (ns *StaffShiftTypes) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
