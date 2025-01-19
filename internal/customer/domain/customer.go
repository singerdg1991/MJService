package domain

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/*
 * @apiDefine: CustomerUserLanguageSkill
 */
type CustomerUserLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:english"`
}

/*
 * @apiDefine: CustomerUser
 */
type CustomerUser struct {
	ID           int64      `json:"id" openapi:"example:1"`
	FirstName    string     `json:"firstName" openapi:"example:saeed"`
	LastName     string     `json:"lastName" openapi:"example:ghanbari"`
	AvatarUrl    string     `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
	Gender       string     `json:"gender" openapi:"example:male"`
	Email        string     `json:"email" openapi:"example:email@yahoo.com"`
	Phone        string     `json:"phone" openapi:"example:09123456789"`
	BirthDate    *time.Time `json:"birthDate" openapi:"example:1990-01-01"`
	NationalCode string     `json:"nationalCode" openapi:"example:1234567890"`
}

/*
 * @apiDefine: CustomerResponsibleNurse
 */
type CustomerResponsibleNurse struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
}

/*
 * @apiDefine: CustomerMotherLang
 */
type CustomerMotherLang struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:english"`
}

/*
 * @apiDefine: CustomersRelative
 */
type CustomersRelative struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName"`
	LastName  string `json:"lastName" openapi:"example:lastName"`
	Relation  string `json:"relation" openapi:"example:father"`
}

/*
 * @apiDefine: Customer
 */
type Customer struct {
	ID                                           int64                     `json:"id" openapi:"example:1"`
	MongoID                                      primitive.ObjectID        `bson:"_id,omitempty" json:"-" openapi:"example:5f7b5f5b9b9b9b9b9b9b9b9b"`
	UserID                                       *int64                    `json:"-" openapi:"example:1"`
	ResponsibleNurseID                           *int64                    `json:"-" openapi:"example:1"`
	RelativeIDs                                  []int64                   `json:"-" openapi:"example:1"`
	DiagnoseIDs                                  []int64                   `json:"-" openapi:"example:1"`
	CreditDetailIDs                              []int64                   `json:"-" openapi:"example:1"`
	ServiceIDs                                   []int64                   `json:"-" openapi:"example:1"`
	AbsenceIDs                                   []int64                   `json:"-" openapi:"example:1"`
	MotherLangIDs                                []int64                   `json:"-" openapi:"example:1"`
	User                                         *CustomerUser             `json:"user" openapi:"ignored"`
	Sections                                     []CustomerSection         `json:"sections" openapi:"ignored"`
	ResponsibleNurse                             *CustomerResponsibleNurse `json:"responsibleStaff" openapi:"ignored"`
	Relatives                                    []CustomersRelative       `json:"relatives" openapi:"ignored"`
	Diagnoses                                    []CustomerDiagnose        `json:"diagnoses" openapi:"ignored"`
	Addresses                                    []CustomerAddress         `json:"addresses" openapi:"ignored"`
	CreditDetails                                []CustomerCreditDetail    `json:"creditDetails" openapi:"ignored"`
	Services                                     []CustomerServices        `json:"services" openapi:"ignored"`
	Absences                                     []CustomerAbsence         `json:"absences" openapi:"ignored"`
	MotherLangs                                  []CustomerMotherLang      `json:"motherLangs" openapi:"ignored"`
	NurseGenderWish                              *string                   `json:"staffGenderWish" openapi:"example:male"`
	Status                                       string                    `json:"status" openapi:"example:active"`
	StatusDate                                   time.Time                 `json:"statusDate" openapi:"example:2021-01-01T00:00:00Z" bson:"statusDate"`
	ParkingInfo                                  *string                   `json:"parkingInfo" openapi:"example:parkingInfo"`
	Limitations                                  []CustomerLimitation      `json:"limitations" openapi:"ignored"`
	ExtraExplanation                             *string                   `json:"extraExplanation" openapi:"example:extraExplanation"`
	HasLimitingTheRightToSelfDetermination       bool                      `json:"hasLimitingTheRightToSelfDetermination" openapi:"example:true"`
	LimitingTheRightToSelfDeterminationStartDate *time.Time                `json:"limitingTheRightToSelfDeterminationStartDate" openapi:"example:2021-01-01T00:00:00Z"`
	LimitingTheRightToSelfDeterminationEndDate   *time.Time                `json:"limitingTheRightToSelfDeterminationEndDate" openapi:"example:2021-01-01T00:00:00Z"`
	MobilityContract                             *string                   `json:"mobilityContract" openapi:"example:mobilityContract"`
	KeyNo                                        *string                   `json:"keyNo" openapi:"example:keyNo"`
	PaymentMethod                                *string                   `json:"paymentMethod" openapi:"example:own"`
	CreatedAt                                    time.Time                 `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                                    time.Time                 `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                                    *time.Time                `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewCustomer() *Customer {
	return &Customer{}
}

func (u *Customer) TableName() string {
	return "customers"
}

func (u *Customer) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *Customer) ToMap() (map[string]interface{}, error) {
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
