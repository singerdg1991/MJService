package domain

/*
 * @apiDefine: ShiftCountPerDay
 */
type ShiftCountPerDay struct {
	DayOfWeek  string `json:"dayOfWeek" openapi:"example:Monday"`
	ShiftCount int    `json:"shiftCount" openapi:"example:25"`
	DayOrder   int    `json:"dayOrder" openapi:"example:1"` // 1 for Monday, 2 for Tuesday, etc.
}
