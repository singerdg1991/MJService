package domain

import (
	"encoding/json"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/*
 * @apiDefine: StaffLanguageSkills
 */
type StaffLanguageSkills struct {
	Name string `json:"name" openapi:"example:car"`
}

/*
 * @apiDefine: StaffVehicleType
 */
type StaffVehicleType struct {
	Name string `json:"name" openapi:"example:car"`
}

/*
 * @apiDefine: StaffTeamRes
 */
type StaffTeamRes struct {
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffAbilityRes
 */
type StaffAbilityRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffUserLanguageSkill
 */
type StaffUserLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:english"`
}

/*
 * @apiDefine: StaffSectionRes
 */
type StaffSectionRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:english"`
}

/*
 * @apiDefine: StaffLicensesResLicense
 */
type StaffLicensesResLicense struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:license"`
}

/*
 * @apiDefine: StaffLicensesRes
 */
type StaffLicensesRes struct {
	ID          uint                    `json:"id" openapi:"example:1"`
	StaffID     uint                    `json:"staffId" openapi:"ignored"`
	LicenseID   uint                    `json:"licenseId" openapi:"ignored"`
	License     StaffLicensesResLicense `json:"license" openapi:"$ref:StaffLicensesRes"`
	ExpireDate  *time.Time              `json:"expire_date" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments []*types.UploadMetadata `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
}

/*
 * @apiDefine: StaffPaymentTypeRes
 */
type StaffPaymentTypeRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:license"`
}

/*
 * @apiDefine: StaffShiftTypeRes
 */
type StaffShiftTypeRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:shiftType"`
}

/*
 * @apiDefine: StaffContractTypeRes
 */
type StaffContractTypeRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:contractType"`
}

/*
 * @apiDefine: StaffUserRole
 */
type StaffUserRole struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:roleName"`
}

/*
 * @apiDefine: StaffUser
 */
type StaffUser struct {
	ID                 uint                     `json:"id" openapi:"example:1"`
	Username           string                   `json:"username" openapi:"example:sgh370"`
	FirstName          string                   `json:"firstName" openapi:"example:saeed"`
	LastName           string                   `json:"lastName" openapi:"example:ghanbari"`
	Email              string                   `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone              string                   `json:"phone" openapi:"example:09123456789"`
	Roles              []StaffUserRole          `json:"roles" openapi:"$ref:StaffUserRole;type:array;"`
	AvatarUrl          string                   `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
	WorkPhoneNumber    string                   `json:"workPhoneNumber" openapi:"example:09123456789"`
	Gender             string                   `json:"gender" openapi:"example:male"`
	AccountNumber      string                   `json:"accountNumber" openapi:"example:1234567890"`
	Telephone          string                   `json:"telephone" openapi:"example:09123456789"`
	RegistrationNumber string                   `json:"registrationNumber" openapi:"example:1234567890"`
	LanguageSkills     []StaffUserLanguageSkill `json:"languageSkills" openapi:"ignored"`
	NationalCode       string                   `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate          *time.Time               `json:"birthDate" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffAddress
 */
type StaffAddress struct {
	ID      uint   `json:"id" openapi:"example:1"`
	Name    string `json:"name" openapi:"example:Home"`
	City    string `json:"city" openapi:"example:tehran"`
	ZipCode string `json:"zipCode" openapi:"example:1234567890"`
	State   string `json:"state" openapi:"example:tehran"`
}

/*
 * @apiDefine: StaffAbsencesQueryResStatusBy
 */
type StaffAbsencesQueryResStatusBy struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:saeed"`
	LastName  string `json:"lastName" openapi:"example:ghanbari"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
}

/*
 * @apiDefine: StaffAbsencesQueryRes
 */
type StaffAbsencesQueryRes struct {
	ID          uint                           `json:"id" openapi:"example:1"`
	StaffID     uint                           `json:"staffId" openapi:"ignored"`
	StartDate   time.Time                      `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate     *time.Time                     `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Reason      *string                        `json:"reason" openapi:"example:reason"`
	Attachments []*types.UploadMetadata        `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Status      *string                        `json:"status" openapi:"example:status"`
	StatusBy    *StaffAbsencesQueryResStatusBy `json:"statusBy" openapi:"$ref:StaffAbsencesQueryResStatusBy"`
	StatusAt    *time.Time                     `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffAbsenceRes
 */
type StaffAbsenceRes struct {
	ID          uint                           `json:"id" openapi:"example:1"`
	StaffID     uint                           `json:"staffId" openapi:"ignored"`
	StartDate   time.Time                      `json:"start_date" openapi:"example:2021-01-01T00:00:00Z"`
	EndDate     *time.Time                     `json:"end_date" openapi:"example:2021-01-01T00:00:00Z"`
	Reason      *string                        `json:"reason" openapi:"example:reason"`
	Attachments []*types.UploadMetadata        `json:"attachments" openapi:"$ref:UploadMetadata;example:[];type:array;required"`
	Status      *string                        `json:"status" openapi:"example:status"`
	StatusBy    *StaffAbsencesQueryResStatusBy `json:"statusBy" openapi:"$ref:StaffAbsencesQueryResStatusBy"`
	StatusAt    *time.Time                     `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffTypesRes
 */
type StaffTypesRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:contract"`
}

/*
 * @apiDefine: Staff
 */
type Staff struct {
	ID                        uint                            `json:"id" openapi:"example:1"`
	MongoID                   primitive.ObjectID              `bson:"_id,omitempty" json:"-" openapi:"example:5f7b5f5b9b9b9b9b9b9b9b9b"`
	UserID                    uint                            `json:"-" openapi:"example:1"`
	PaymentTypeID             *uint                           `json:"-" openapi:"example:1"`
	Addresses                 []StaffAddress                  `json:"addresses" openapi:"ignored"`
	User                      StaffUser                       `json:"user" openapi:"ignored"`
	Sections                  []StaffSectionRes               `json:"sections" openapi:"ignored"`
	Licenses                  []StaffLicensesRes              `json:"licenses" openapi:"ignored"`
	PaymentType               *StaffPaymentTypeRes            `json:"paymentType" openapi:"ignored"`
	ShiftTypes                []StaffShiftTypeRes             `json:"availableShifts" openapi:"ignored"`
	ContractTypes             []StaffContractTypeRes          `json:"contractTypes" openapi:"ignored"`
	Absences                  []StaffAbsenceRes               `json:"absences" openapi:"ignored"`
	StaffTypes                []StaffTypesRes                 `json:"staffTypes" openapi:"ignored"`
	Limitations               []sharedmodels.SharedLimitation `json:"limitations" openapi:"$ref:SharedLimitation;type:array;"`
	CertificateCode           *string                         `json:"certificateCode" openapi:"example:1234567890"`
	JobTitle                  *string                         `json:"jobTitle" openapi:"example:jobTitle"`
	Grace                     int                             `json:"grace" openapi:"example:1"`
	Warning                   int                             `json:"warning" openapi:"example:1"`
	Attention                 int                             `json:"attention" openapi:"example:1"`
	Progress                  int                             `json:"progress" openapi:"example:1"`
	ExperienceAmount          *uint                           `json:"experienceAmount" openapi:"example:1"`
	ExperienceAmountUnit      *string                         `json:"experienceAmountUnit" openapi:"example:month"`
	IsSubcontractor           bool                            `json:"isSubcontractor" openapi:"example:true"`
	CompanyRegistrationNumber *string                         `json:"companyRegistrationNumber" openapi:"example:1234567890"`
	OrganizationNumber        *string                         `json:"organizationNumber" openapi:"example:1234567890"`
	PercentLengthInContract   *uint                           `json:"percentLengthInContract" openapi:"example:100"`
	HourLengthInContract      *uint                           `json:"hourLengthInContract" openapi:"example:40"`
	Salary                    *uint                           `json:"salary" openapi:"example:1000000"`
	VehicleTypes              interface{}                     `json:"vehicleTypes" openapi:"example:[\"car\",\"bicycle\",\"public_transportation\"];type:array;required;"`
	VehicleLicenseTypes       interface{}                     `json:"vehicleLicenseTypes" openapi:"example:[\"automatic\",\"manual\"];type:array;required;"`
	TrialTime                 *time.Time                      `json:"trial_time" openapi:"example:2021-01-01T00:00:00Z"`
	Attachments               []*types.UploadMetadata         `json:"attachments" openapi:"$ref:UploadMetadata;type:array;required"`
	JoinedAt                  *time.Time                      `json:"joined_at" openapi:"example:2021-01-01T00:00:00Z"`
	ContractStartedAt         *time.Time                      `json:"contract_started_at" openapi:"example:2021-01-01T00:00:00Z"`
	ContractExpiresAt         *time.Time                      `json:"contract_expires_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt                 time.Time                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt                 time.Time                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt                 *time.Time                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

func NewStaff() *Staff {
	return &Staff{}
}

func (u *Staff) TableName() string {
	return "staffs"
}

func (u *Staff) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *Staff) ToMap() (map[string]interface{}, error) {
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
