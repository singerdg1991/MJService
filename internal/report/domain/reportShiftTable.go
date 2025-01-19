package domain

import (
	"encoding/json"
	"time"

	cycleDomain "github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: ReportShiftTableStaff
 */
type ReportShiftTableStaff struct {
	ID        string `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: ReportShiftTableCustomer
 */
type ReportShiftTableCustomer struct {
	ID        string `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:John;required"`
	LastName  string `json:"lastName" openapi:"example:Doe;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com"`
}

/*
 * @apiDefine: ReportShiftTable
 */
type ReportShiftTable struct {
	ID            uint                          `json:"id" openapi:"example:1"`
	ExchangeKey   string                        `json:"exchangeKey" openapi:"example:dfhdsjrtwerwrwfgjgfrt"`
	CycleID       uint                          `json:"cycleId" openapi:"example:1"`
	StaffTypeIDs  []uint                        `json:"staffTypeIds" openapi:"example:[1,2,3]"`
	StaffTypes    []*cycleDomain.CycleStaffType `json:"staffTypes" openapi:"$ref:CycleStaffType;type:array"`
	ShiftName     string                        `json:"shiftName" openapi:"example:morning"`
	VehicleType   *string                       `json:"vehicleType" openapi:"example:own"`
	StartLocation *string                       `json:"startLocation" openapi:"example:office"`
	DateTime      time.Time                     `json:"dateTime" openapi:"example:2021-08-02"`
	Status        string                        `json:"status" openapi:"example:not-started"`
	CreatedAt     time.Time                     `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time                     `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt     *time.Time                    `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewReportShift() *ReportShiftTable {
	return &ReportShiftTable{}
}

func (u *ReportShiftTable) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *ReportShiftTable) ToMap() (map[string]interface{}, error) {
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
