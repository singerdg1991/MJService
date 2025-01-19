package domain

import "time"

/*
 * @apiDefine: ShiftSchedulingChartData
 */
type ShiftSchedulingChartData struct {
	StaffID     uint      `json:"staffId" openapi:"example:1"`
	StaffName   string    `json:"staffName" openapi:"example:John Doe"`
	ShiftName   string    `json:"shiftName" openapi:"example:morning"`
	StartHour   time.Time `json:"startHour" openapi:"example:2021-01-01T09:00:00Z"`
	EndHour     time.Time `json:"endHour" openapi:"example:2021-01-01T17:00:00Z"`
	Date        time.Time `json:"date" openapi:"example:2021-01-01T00:00:00Z"`
	Status      string    `json:"status" openapi:"example:assigned"`
	IsUnplanned bool      `json:"isUnplanned" openapi:"example:false"`
}
