package models

import (
	customerDomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/report/domain"
)

/*
 * @apiDefine: ReportsQueryCustomersTableResponseDataItem
 */
type ReportsQueryCustomersTableResponseDataItem struct {
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
	StatusDate                                   string                                   `json:"statusDate" openapi:"example:2021-01-01T00:00:00Z" bson:"statusDate"`
	ParkingInfo                                  *string                                  `json:"parkingInfo" openapi:"example:parkingInfo"`
	Limitations                                  []customerDomain.CustomerLimitation      `json:"limitations" openapi:"ignored"`
	ExtraExplanation                             *string                                  `json:"extraExplanation" openapi:"example:extraExplanation"`
	HasLimitingTheRightToSelfDetermination       bool                                     `json:"hasLimitingTheRightToSelfDetermination" openapi:"example:true"`
	LimitingTheRightToSelfDeterminationStartDate *string                                  `json:"limitingTheRightToSelfDeterminationStartDate" openapi:"example:2021-01-01T00:00:00Z"`
	LimitingTheRightToSelfDeterminationEndDate   *string                                  `json:"limitingTheRightToSelfDeterminationEndDate" openapi:"example:2021-01-01T00:00:00Z"`
	MobilityContract                             *string                                  `json:"mobilityContract" openapi:"example:mobilityContract"`
	KeyNo                                        *string                                  `json:"keyNo" openapi:"example:keyNo"`
	PaymentMethod                                *string                                  `json:"paymentMethod" openapi:"example:own"`
	CreatedAt                                    string                                   `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                                    string                                   `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                                    *string                                  `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: ReportsQueryCustomersTableResponseData
 */
type ReportsQueryCustomersTableResponseData struct {
	Limit      int                                          `json:"limit" openapi:"example:10"`
	Offset     int                                          `json:"offset" openapi:"example:0"`
	Page       int                                          `json:"page" openapi:"example:1"`
	TotalRows  int                                          `json:"totalRows" openapi:"example:1"`
	TotalPages int                                          `json:"totalPages" openapi:"example:1"`
	Items      []ReportsQueryCustomersTableResponseDataItem `json:"items" openapi:"$ref:ReportsQueryCustomersTableResponseDataItem;type:array"`
}

/*
 * @apiDefine: ReportsQueryCustomersTableResponse
 */
type ReportsQueryCustomersTableResponse struct {
	StatusCode int                                    `json:"statusCode" openapi:"example:200"`
	Data       ReportsQueryCustomersTableResponseData `json:"data" openapi:"$ref:ReportsQueryCustomersTableResponseData"`
}

/*
 * @apiDefine: ReportsQueryCustomersTableNotFoundResponse
 */
type ReportsQueryCustomersTableNotFoundResponse struct {
	ReportsQueryCustomersTable []domain.ReportCustomerTable `json:"reportsQueryCustomersTable" openapi:"$ref:ReportCustomerTable;type:array"`
}
