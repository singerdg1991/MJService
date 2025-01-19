package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
)

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataUser
 */
type StaffsCreateOrUpdateContractResponseDataAddress struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataUserLanguageSkill
 */
type StaffsCreateOrUpdateContractResponseDataUserLanguageSkill struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataUser
 */
type StaffsCreateOrUpdateContractResponseDataUser struct {
	ID                 int                                                         `json:"id" openapi:"example:1"`
	Username           string                                                      `json:"username" openapi:"example:test"`
	FirstName          string                                                      `json:"firstName" openapi:"example:test"`
	LastName           string                                                      `json:"lastName" openapi:"example:test"`
	Email              string                                                      `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone              string                                                      `json:"phone" openapi:"example:123456789"`
	Roles              []domain.StaffUserRole                                      `json:"roles" openapi:"$ref:StaffUserRole;type:array;"`
	AvatarUrl          string                                                      `json:"avatarUrl" openapi:"example:https://www.google.com"`
	WorkPhoneNumber    string                                                      `json:"workPhoneNumber" openapi:"example:123456789"`
	Gender             string                                                      `json:"gender" openapi:"example:male"`
	AccountNumber      string                                                      `json:"accountNumber" openapi:"example:123456789"`
	Telephone          string                                                      `json:"telephone" openapi:"example:123456789"`
	RegistrationNumber string                                                      `json:"registrationNumber" openapi:"example:123456789"`
	LanguageSkills     []StaffsCreateOrUpdateContractResponseDataUserLanguageSkill `json:"languageSkills" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataUserLanguageSkill;type:array"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataSection
 */
type StaffsCreateOrUpdateContractResponseDataSection struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataPermission
 */
type StaffsCreateOrUpdateContractResponseDataPermission struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataPaymentType
 */
type StaffsCreateOrUpdateContractResponseDataPaymentType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataAvailableShift
 */
type StaffsCreateOrUpdateContractResponseDataAvailableShift struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataContractType
 */
type StaffsCreateOrUpdateContractResponseDataContractType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataAbsence
 */
type StaffsCreateOrUpdateContractResponseDataAbsence struct {
	ID        int     `json:"id" openapi:"example:1"`
	StartDate string  `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate   *string `json:"end_date" openapi:"example:2020-01-01T00:00:00Z;nullable"`
	Reason    *string `json:"reason" openapi:"example:reason;nullable"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseDataStaffType
 */
type StaffsCreateOrUpdateContractResponseDataStaffType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponseData
 */
type StaffsCreateOrUpdateContractResponseData struct {
	ID                        int                                                      `json:"id" openapi:"example:1"`
	Addresses                 []StaffsCreateOrUpdateContractResponseDataAddress        `json:"addresses" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataUser;type:array"`
	User                      StaffsCreateOrUpdateContractResponseDataUser             `json:"user" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataUser"`
	Sections                  []StaffsCreateOrUpdateContractResponseDataSection        `json:"sections" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataSection;type:array"`
	Permissions               []StaffsCreateOrUpdateContractResponseDataPermission     `json:"permissions" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataPermission;type:array"`
	PaymentType               StaffsCreateOrUpdateContractResponseDataPaymentType      `json:"paymentType" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataPaymentType"`
	AvailableShifts           []StaffsCreateOrUpdateContractResponseDataAvailableShift `json:"availableShifts" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataAvailableShift;type:array"`
	ContractTypes             []StaffsCreateOrUpdateContractResponseDataContractType   `json:"contractTypes" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataContractType;type:array"`
	Absences                  []StaffsCreateOrUpdateContractResponseDataAbsence        `json:"absences" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataAbsence;type:array"`
	Limitations               string                                                   `json:"limitations" openapi:"example:lorem ipsum"`
	JobTitle                  string                                                   `json:"jobTitle" openapi:"example:lorem ipsum"`
	CertificateCode           string                                                   `json:"certificateCode" openapi:"example:1234567890"`
	Grace                     int                                                      `json:"grace" openapi:"example:1"`
	Warning                   int                                                      `json:"warning" openapi:"example:1"`
	Attention                 int                                                      `json:"attention" openapi:"example:1"`
	Progress                  int                                                      `json:"progress" openapi:"example:1"`
	ExperienceAmount          int                                                      `json:"experienceAmount" openapi:"example:1"`
	ExperienceAmountUnit      string                                                   `json:"experienceAmountUnit" openapi:"example:1"`
	IsSubcontractor           bool                                                     `json:"isSubcontractor" openapi:"example:true"`
	CompanyRegistrationNumber string                                                   `json:"companyRegistrationNumber" openapi:"example:1"`
	PercentLengthInContract   int                                                      `json:"percentLengthInContract" openapi:"example:1"`
	HourLengthInContract      int                                                      `json:"hourLengthInContract" openapi:"example:1"`
	Salary                    int                                                      `json:"salary" openapi:"example:1"`
	StaffTypes                []StaffsCreateOrUpdateContractResponseDataStaffType      `json:"staffTypes" openapi:"$ref:StaffsCreateOrUpdateContractResponseDataStaffType;type:array"`
	TrialTime                 int                                                      `json:"trial_time" openapi:"example:2022-01-01T00:00:00Z"`
	Attachments               []types.UploadMetadata                                   `json:"attachments" openapi:"$ref:UploadMetadata;type:array;required"`
	JoinedAt                  int                                                      `json:"joined_at" openapi:"example:2022-01-01T00:00:00Z"`
	ContractStartedAt         int                                                      `json:"contract_started_at" openapi:"example:2022-01-01T00:00:00Z"`
	ContractExpiresAt         int                                                      `json:"contract_expires_at" openapi:"example:2022-01-01T00:00:00Z"`
	CreatedAt                 int                                                      `json:"created_at" openapi:"example:2022-01-01T00:00:00Z"`
	UpdatedAt                 int                                                      `json:"updated_at" openapi:"example:2022-01-01T00:00:00Z"`
	DeletedAt                 int                                                      `json:"deleted_at" openapi:"example:2022-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractResponse
 */
type StaffsCreateOrUpdateContractResponse struct {
	StatusCode int                                      `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreateOrUpdateContractResponseData `json:"data" openapi:"$ref:StaffsCreateOrUpdateContractResponseData"`
}
