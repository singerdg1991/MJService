package domain

import (
	"encoding/json"
	"time"
)

type StaffSection struct {
	ID        int64      `json:"id"`
	StaffID   int64      `json:"staffId"`
	SectionID int64      `json:"sectionId"`
	CreatedAt time.Time  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffSection() *StaffSection {
	return &StaffSection{}
}

func (ns *StaffSection) TableName() string {
	return "staffSections"
}

func (ns *StaffSection) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
