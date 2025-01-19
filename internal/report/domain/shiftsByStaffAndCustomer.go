package domain

/*
 * @apiDefine: ShiftsByStaffAndCustomer
 */
type ShiftsByStaffAndCustomer struct {
	StaffId      int64  `json:"staffId" openapi:"example:1"`
	StaffName    string `json:"staffName" openapi:"example:John Doe"`
	CustomerId   int64  `json:"customerId" openapi:"example:1"`
	CustomerName string `json:"customerName" openapi:"example:Customer A"`
	ShiftCount   int    `json:"shiftCount" openapi:"example:30"`
}
