package domain

import (
	"encoding/json"
	"time"
)

/*
 * @apiDefine: ReportArrangementTableStaff
 */
type ReportArrangementTableStaff struct {
	ID        uint   `json:"id" openapi:"example:1"`
	UserID    uint   `json:"userId" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://example.com/avatar.jpg"`
}

/*
 * @apiDefine: ReportArrangementTableShift
 */
type ReportArrangementTableShift struct {
	ID          uint      `json:"id" openapi:"example:1"`
	ShiftName   string    `json:"shiftName" openapi:"example:morning"`
	StartHour   time.Time `json:"startHour" openapi:"example:00:00"`
	EndHour     time.Time `json:"endHour" openapi:"example:00:00"`
	IsUnplanned bool      `json:"isUnplanned" openapi:"example:false"`
	DateTime    time.Time `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: ReportArrangementTable
 */
type ReportArrangementTable struct {
	ID                      uint                       `json:"id" openapi:"example:1"`
	CycleID                 uint                       `json:"cycleId" openapi:"example:1"`
	Staff                   *ReportArrangementTableStaff `json:"staff" openapi:"$ref:ReportArrangementTableStaff"`
	Shift                   *ReportArrangementTableShift `json:"shift" openapi:"$ref:ReportArrangementTableShift"`
	Status                  string                     `json:"status" openapi:"example:not-started"`
	PrevStatus              string                     `json:"prevStatus" openapi:"example:not-started"`
	StartKilometer          *string                    `json:"startKilometer" openapi:"example:0"`
	ReasonOfTheCancellation *string                    `json:"reasonOfTheCancellation" openapi:"example:reason"`
	ReasonOfTheReactivation *string                    `json:"reasonOfTheReactivation" openapi:"example:reason"`
	ReasonOfTheResume       *string                    `json:"reasonOfTheResume" openapi:"example:reason"`
	ReasonOfThePause        *string                    `json:"reasonOfThePause" openapi:"example:reason"`
	IsUnplanned             bool                       `json:"isUnplanned" openapi:"example:false"`
	DateTime                time.Time                  `json:"datetime" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt               time.Time                  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt               time.Time                  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt               *time.Time                 `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
	StartedAt               *time.Time                 `json:"started_at" openapi:"example:2021-01-01T00:00:00Z"`
	EndedAt                 *time.Time                 `json:"ended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CancelledAt             *time.Time                 `json:"cancelled_at" openapi:"example:2021-01-01T00:00:00Z"`
	DelayedAt               *time.Time                 `json:"delayed_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewReportArrangement() *ReportArrangementTable {
	return &ReportArrangementTable{}
}

func (u *ReportArrangementTable) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *ReportArrangementTable) ToMap() (map[string]interface{}, error) {
	jsonString, err := u.ToJson()
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
