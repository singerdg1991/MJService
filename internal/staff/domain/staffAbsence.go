package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"time"
)

/*
 * @apiDefine: StaffAbsencesStatusBy
 */
type StaffAbsencesStatusBy struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:saeed"`
	LastName  string `json:"lastName" openapi:"example:ghanbari"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
}

type StaffAbsences struct {
	ID          int64                   `json:"id"`
	StaffID     int64                   `json:"staffId"`
	StartDate   time.Time               `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate     *time.Time              `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Reason      *string                 `json:"reason" openapi:"example:reason"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Status      *string                 `json:"status" openapi:"example:status"`
	StatusBy    *StaffAbsencesStatusBy  `json:"statusBy" openapi:"$ref:StaffAbsencesStatusBy"`
	StatusAt    *time.Time              `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt   time.Time               `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time               `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *time.Time              `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaffAbsences() *StaffAbsences {
	return &StaffAbsences{}
}

func (ns *StaffAbsences) TableName() string {
	return "staffAbsences"
}

func (ns *StaffAbsences) ToJson() (string, error) {
	b, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
