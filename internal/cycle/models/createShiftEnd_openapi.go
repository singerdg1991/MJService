package models

/*
 * @apiDefine: CyclesCreateShiftEndResponseData
 */
type CyclesCreateShiftEndResponseData struct {
	ID            uint    `json:"id" openapi:"example:1"`
	ExchangeKey   string  `json:"exchangeKey" openapi:"example:dfhdsjrtwerwrwfgjgfrt"`
	CycleID       uint    `json:"cycleId" openapi:"example:1"`
	StaffTypeIDs  []uint  `json:"staffTypeIds" openapi:"example:[1,2,3]"`
	ShiftName     string  `json:"shiftName" openapi:"example:morning"`
	VehicleType   *string `json:"vehicleType" openapi:"example:own"`
	StartLocation *string `json:"startLocation" openapi:"example:office"`
	DateTime      string  `json:"dateTime" openapi:"example:2021-08-02"`
	Status        string  `json:"status" openapi:"example:not-started"`
	CreatedAt     string  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     string  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *string `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CyclesCreateShiftEndResponse
 */
type CyclesCreateShiftEndResponse struct {
	StatusCode int                              `json:"statusCode" openapi:"example:200"`
	Data       CyclesCreateShiftEndResponseData `json:"data" openapi:"$ref:CyclesCreateShiftEndResponseData"`
}
