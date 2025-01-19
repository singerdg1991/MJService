package domain

import (
	"encoding/json"
	"time"

	customerDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: ReportCustomerTable
 */
type ReportCustomerTable struct {
	ID                                           int64                                    `json:"id" openapi:"example:1"`
	UserID                                       *int64                                   `json:"-" openapi:"example:1"`
	ResponsibleNurseID                           *int64                                   `json:"-" openapi:"example:1"`
	RelativeIDs                                  []int64                                  `json:"-" openapi:"example:1"`
	DiagnoseIDs                                  []int64                                  `json:"-" openapi:"example:1"`
	CreditDetailIDs                              []int64                                  `json:"-" openapi:"example:1"`
	ServiceIDs                                   []int64                                  `json:"-" openapi:"example:1"`
	AbsenceIDs                                   []int64                                  `json:"-" openapi:"example:1"`
	MotherLangIDs                                []int64                                  `json:"-" openapi:"example:1"`
	User                                         *customerDomain.CustomerUser             `json:"user" openapi:"ignored"`
	Sections                                     []customerDomain.CustomerSection         `json:"sections" openapi:"ignored"`
	ResponsibleNurse                             *customerDomain.CustomerResponsibleNurse `json:"responsibleStaff" openapi:"ignored"`
	Relatives                                    []customerDomain.CustomersRelative       `json:"relatives" openapi:"ignored"`
	Diagnoses                                    []customerDomain.CustomerDiagnose        `json:"diagnoses" openapi:"ignored"`
	Addresses                                    []customerDomain.CustomerAddress         `json:"addresses" openapi:"ignored"`
	CreditDetails                                []customerDomain.CustomerCreditDetail    `json:"creditDetails" openapi:"ignored"`
	Services                                     []customerDomain.CustomerServices        `json:"services" openapi:"ignored"`
	Absences                                     []customerDomain.CustomerAbsence         `json:"absences" openapi:"ignored"`
	MotherLangs                                  []customerDomain.CustomerMotherLang      `json:"motherLangs" openapi:"ignored"`
	NurseGenderWish                              *string                                  `json:"staffGenderWish" openapi:"example:male"`
	Status                                       string                                   `json:"status" openapi:"example:active"`
	StatusDate                                   time.Time                                `json:"statusDate" openapi:"example:2021-01-01T00:00:00Z" bson:"statusDate"`
	ParkingInfo                                  *string                                  `json:"parkingInfo" openapi:"example:parkingInfo"`
	Limitations                                  []customerDomain.CustomerLimitation      `json:"limitations" openapi:"ignored"`
	ExtraExplanation                             *string                                  `json:"extraExplanation" openapi:"example:extraExplanation"`
	HasLimitingTheRightToSelfDetermination       bool                                     `json:"hasLimitingTheRightToSelfDetermination" openapi:"example:true"`
	LimitingTheRightToSelfDeterminationStartDate *time.Time                               `json:"limitingTheRightToSelfDeterminationStartDate" openapi:"example:2021-01-01T00:00:00Z"`
	LimitingTheRightToSelfDeterminationEndDate   *time.Time                               `json:"limitingTheRightToSelfDeterminationEndDate" openapi:"example:2021-01-01T00:00:00Z"`
	MobilityContract                             *string                                  `json:"mobilityContract" openapi:"example:mobilityContract"`
	KeyNo                                        *string                                  `json:"keyNo" openapi:"example:keyNo"`
	PaymentMethod                                *string                                  `json:"paymentMethod" openapi:"example:own"`
	CreatedAt                                    time.Time                                `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                                    time.Time                                `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                                    *time.Time                               `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewReportCustomer() *ReportCustomerTable {
	return &ReportCustomerTable{}
}

func (u *ReportCustomerTable) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *ReportCustomerTable) ToMap() (map[string]interface{}, error) {
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
