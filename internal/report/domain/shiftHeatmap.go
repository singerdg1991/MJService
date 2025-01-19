package domain

// ShiftHeatmap represents shift distribution across days and hours
type ShiftHeatmap struct {
	DayOfWeek  string `json:"dayOfWeek" openapi:"example:Monday"` // Monday, Tuesday, etc.
	DayOrder   int    `json:"dayOrder" openapi:"example:1"`       // 1-7 for sorting
	Hour       int    `json:"hour" openapi:"example:0"`           // 0-23
	ShiftCount int    `json:"shiftCount" openapi:"example:0"`     // Number of shifts in this time slot
}
