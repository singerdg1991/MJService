package domain

/*
 * @apiDefine: ShiftDistributionByCustomer
 */
type ShiftDistributionByCustomer struct {
	CustomerID     uint    `json:"customerId" openapi:"example:1"`
	CustomerName   string  `json:"customerName" openapi:"example:Customer A"`
	NumberOfShifts int     `json:"numberOfShifts" openapi:"example:30"`
	Percentage     float64 `json:"percentage" openapi:"example:25.5"`
}
