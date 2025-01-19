package domain

/*
 * @apiDefine: ShiftDurationAnalysis
 */
type ShiftDurationAnalysis struct {
	CustomerID       uint    `json:"customerId" openapi:"example:1"`
	CustomerName     string  `json:"customerName" openapi:"example:Customer A"`
	DurationHours    float64 `json:"durationHours" openapi:"example:6"`
	NumberOfShifts   int     `json:"numberOfShifts" openapi:"example:30"`
}
