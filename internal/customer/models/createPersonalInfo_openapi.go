package models

import (
	"github.com/hoitek/Maja-Service/internal/customer/domain"
)

/*
 * @apiDefine: CustomersCreatePersonalInfoResponseDataUser
 */
type CustomersCreatePersonalInfoResponseDataUser struct {
	ID           int64  `json:"id" openapi:"example:1"`
	FirstName    string `json:"firstName" openapi:"example:saeed"`
	LastName     string `json:"lastName" openapi:"example:ghanbari"`
	AvatarUrl    string `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
	Gender       string `json:"gender" openapi:"example:male"`
	Email        string `json:"email" openapi:"example:email@yahoo.com"`
	Phone        string `json:"phone" openapi:"example:09123456789"`
	BirthDate    string `json:"birthDate" openapi:"example:1990-01-01"`
	NationalCode string `json:"nationalCode" openapi:"example:1234567890"`
}

/*
 * @apiDefine: CustomersCreatePersonalInfoResponseData
 */
type CustomersCreatePersonalInfoResponseData struct {
	ID                                           int                                         `json:"id" openapi:"example:1"`
	User                                         CustomersCreatePersonalInfoResponseDataUser `json:"user" openapi:"$ref:CustomerUser"`
	Sections                                     []domain.CustomerSection                    `json:"sections" openapi:"$ref:CustomerSection;type:array"`
	ResponsibleNurse                             domain.CustomerResponsibleNurse             `json:"responsibleStaff" openapi:"$ref:CustomerResponsibleNurse"`
	Relatives                                    []domain.CustomersRelative                  `json:"relatives" openapi:"$ref:CustomersRelative;type:array"`
	Diagnoses                                    []domain.CustomerDiagnose                   `json:"diagnoses" openapi:"$ref:CustomerDiagnose;type:array"`
	Addresses                                    []domain.CustomerAddress                    `json:"addresses" openapi:"$ref:CustomerAddress;type:array"`
	CreditDetails                                []domain.CustomerCreditDetail               `json:"creditDetails" openapi:"$ref:CustomerCreditDetail;type:array"`
	Services                                     []domain.CustomerServices                   `json:"services" openapi:"$ref:CustomerServices;type:array"`
	Absences                                     []domain.CustomerAbsence                    `json:"absences" openapi:"$ref:CustomerAbsence;type:array"`
	MotherLangs                                  []domain.CustomerMotherLang                 `json:"motherLangs" openapi:"$ref:CustomerMotherLang;type:array"`
	NurseGenderWish                              string                                      `json:"staffGenderWish" openapi:"example:male"`
	Status                                       string                                      `json:"status" openapi:"example:active"`
	StatusDate                                   string                                      `json:"statusDate" openapi:"example:2020-01-01T00:00:00Z"`
	ParkingInfo                                  string                                      `json:"parkingInfo" openapi:"example:lorem ipsum"`
	Limitations                                  []domain.CustomerLimitation                 `json:"limitations" openapi:"$ref:CustomerLimitation;type:array"`
	ExtraExplanation                             string                                      `json:"extraExplanation" openapi:"example:lorem ipsum"`
	HasLimitingTheRightToSelfDetermination       bool                                        `json:"hasLimitingTheRightToSelfDetermination" openapi:"example:true"`
	LimitingTheRightToSelfDeterminationStartDate string                                      `json:"limitingTheRightToSelfDeterminationStartDate" openapi:"example:2020-01-01T00:00:00Z"`
	LimitingTheRightToSelfDeterminationEndDate   string                                      `json:"limitingTheRightToSelfDeterminationEndDate" openapi:"example:2020-01-01T00:00:00Z"`
	MobilityContract                             string                                      `json:"mobilityContract" openapi:"example:lorem ipsum"`
	KeyNo                                        string                                      `json:"keyNo" openapi:"example:lorem ipsum"`
	PaymentMethod                                string                                      `json:"paymentMethod" openapi:"example:lorem ipsum"`
	CreatedAt                                    string                                      `json:"created_at" openapi:"example:2020-01-01T00:00:00Z"`
	UpdatedAt                                    string                                      `json:"updated_at" openapi:"example:2020-01-01T00:00:00Z"`
	DeletedAt                                    string                                      `json:"deleted_at" openapi:"example:2020-01-01T00:00:00Z"`
}

/*
 * @apiDefine: CustomersCreatePersonalInfoResponse
 */
type CustomersCreatePersonalInfoResponse struct {
	StatusCode int                                     `json:"statusCode" openapi:"example:200"`
	Data       CustomersCreatePersonalInfoResponseData `json:"data" openapi:"$ref:CustomersCreatePersonalInfoResponseData"`
}
