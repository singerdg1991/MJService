package domain

import (
	"encoding/json"
	"time"

	"github.com/hoitek/Maja-Service/internal/cycle/constants"
)

/*
 * @apiDefine: Cycle
 */
type Cycle struct {
	ID                    uint                 `json:"id" openapi:"example:1" bson:"id"`
	SectionID             uint                 `json:"sectionId" openapi:"example:1" bson:"sectionId"`
	Name                  string               `json:"name" openapi:"example:John;required" bson:"name"`
	StartDate             time.Time            `json:"start_date" openapi:"example:2021-01-01;required" bson:"start_date"`
	EndDate               *time.Time           `json:"end_date" openapi:"example:2021-01-01;required" bson:"end_date"`
	PeriodLength          *string              `json:"periodLength" openapi:"example:oneWeek;required" bson:"periodLength"`
	ShiftMorningStartTime string               `json:"shiftMorningStartTime" openapi:"example:08:00;required;" bson:"shiftMorningStartTime"`
	ShiftMorningEndTime   string               `json:"shiftMorningEndTime" openapi:"example:16:00;required;" bson:"shiftMorningEndTime"`
	ShiftEveningStartTime string               `json:"shiftEveningStartTime" openapi:"example:16:00;required;" bson:"shiftEveningStartTime"`
	ShiftEveningEndTime   string               `json:"shiftEveningEndTime" openapi:"example:00:00;required;" bson:"shiftEveningEndTime"`
	ShiftNightStartTime   string               `json:"shiftNightStartTime" openapi:"example:00:00;required;" bson:"shiftNightStartTime"`
	ShiftNightEndTime     string               `json:"shiftNightEndTime" openapi:"example:08:00;required;" bson:"shiftNightEndTime"`
	FreezePeriodDate      time.Time            `json:"freeze_period_date" openapi:"example:2021-01-01;required" bson:"freeze_period_date"`
	WishDays              int                  `json:"wishDays" openapi:"example:1;required" bson:"wishDays"`
	StaffTypes            []CycleStaffType     `json:"staffTypes" openapi:"$ref:CycleStaffType;type:array" bson:"staffTypes"`
	NextStaffTypes        []CycleNextStaffType `json:"nextStaffTypes" openapi:"$ref:CycleNextStaffType;type:array" bson:"cycleNextStaffTypes"`
	Status                string               `json:"status" openapi:"example:active;required" bson:"status"`
	CreatedAt             time.Time            `json:"created_at" openapi:"example:2021-01-01T00:00:00Z" bson:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z" bson:"updated_at"`
	DeletedAt             *time.Time           `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z" bson:"deleted_at"`
}

func NewCycle() *Cycle {
	return &Cycle{}
}

func (c *Cycle) TableName() string {
	return "cycles"
}

func (c *Cycle) ToJson() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Cycle) SetDefaultStatus() string {
	// Set default status
	c.Status = constants.STATUS_ACTIVE

	// Check if cycle is frozen
	if time.Now().After(c.FreezePeriodDate) {
		c.Status = constants.STATUS_FROZEN
	}

	// Check if cycle is expired
	if c.EndDate != nil && time.Now().After(*c.EndDate) {
		c.Status = constants.STATUS_EXPIRED
	}

	return c.Status
}
