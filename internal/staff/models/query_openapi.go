package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"time"
)

/*
 * @apiDefine: StaffsQueryResponseDataUser
 */
type StaffsQueryResponseDataAddress struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataUserLanguageSkill
 */
type StaffsQueryResponseDataUserLanguageSkill struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataUser
 */
type StaffsQueryResponseDataUser struct {
	ID                 int                                        `json:"id" openapi:"example:1"`
	Username           string                                     `json:"username" openapi:"example:test"`
	FirstName          string                                     `json:"firstName" openapi:"example:test"`
	LastName           string                                     `json:"lastName" openapi:"example:test"`
	Email              string                                     `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone              string                                     `json:"phone" openapi:"example:123456789"`
	Roles              []domain.StaffUserRole                     `json:"roles" openapi:"$ref:StaffUserRole;type:array;"`
	AvatarUrl          string                                     `json:"avatarUrl" openapi:"example:https://www.google.com"`
	WorkPhoneNumber    string                                     `json:"workPhoneNumber" openapi:"example:123456789"`
	Gender             string                                     `json:"gender" openapi:"example:male"`
	AccountNumber      string                                     `json:"accountNumber" openapi:"example:123456789"`
	Telephone          string                                     `json:"telephone" openapi:"example:123456789"`
	RegistrationNumber string                                     `json:"registrationNumber" openapi:"example:123456789"`
	LanguageSkills     []StaffsQueryResponseDataUserLanguageSkill `json:"languageSkills" openapi:"$ref:StaffsQueryResponseDataUserLanguageSkill;type:array"`
	NationalCode       string                                     `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate          *string                                    `json:"birthDate" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsQueryResponseDataSection
 */
type StaffsQueryResponseDataSection struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataPermission
 */
type StaffsQueryResponseDataPermission struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataPaymentType
 */
type StaffsQueryResponseDataPaymentType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataAvailableShift
 */
type StaffsQueryResponseDataAvailableShift struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataContractType
 */
type StaffsQueryResponseDataContractType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataAbsence
 */
type StaffsQueryResponseDataAbsence struct {
	ID        int                           `json:"id" openapi:"example:1"`
	StartDate string                        `json:"start_date" openapi:"example:2020-01-01T00:00:00Z"`
	EndDate   *string                       `json:"end_date" openapi:"example:2020-01-01T00:00:00Z;nullable"`
	Reason    *string                       `json:"reason" openapi:"example:reason;nullable"`
	Status    *string                       `json:"status" openapi:"example:status"`
	StatusBy  *domain.StaffAbsencesStatusBy `json:"statusBy" openapi:"$ref:StaffAbsencesStatusBy"`
	StatusAt  *time.Time                    `json:"status_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsQueryResponseDataStaffType
 */
type StaffsQueryResponseDataStaffType struct {
	ID   int    `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test"`
}

/*
 * @apiDefine: StaffsQueryResponseDataItem
 */
type StaffsQueryResponseDataItem struct {
	ID                        int                                     `json:"id" openapi:"example:1"`
	Addresses                 []StaffsQueryResponseDataAddress        `json:"addresses" openapi:"$ref:StaffsQueryResponseDataUser;type:array"`
	User                      StaffsQueryResponseDataUser             `json:"user" openapi:"$ref:StaffsQueryResponseDataUser"`
	Sections                  []StaffsQueryResponseDataSection        `json:"sections" openapi:"$ref:StaffsQueryResponseDataSection;type:array"`
	Permissions               []StaffsQueryResponseDataPermission     `json:"permissions" openapi:"$ref:StaffsQueryResponseDataPermission;type:array"`
	PaymentType               StaffsQueryResponseDataPaymentType      `json:"paymentType" openapi:"$ref:StaffsQueryResponseDataPaymentType"`
	AvailableShifts           []StaffsQueryResponseDataAvailableShift `json:"availableShifts" openapi:"$ref:StaffsQueryResponseDataAvailableShift;type:array"`
	ContractTypes             []StaffsQueryResponseDataContractType   `json:"contractTypes" openapi:"$ref:StaffsQueryResponseDataContractType;type:array"`
	Absences                  []StaffsQueryResponseDataAbsence        `json:"absences" openapi:"$ref:StaffsQueryResponseDataAbsence;type:array"`
	Limitations               []sharedmodels.SharedLimitation         `json:"limitations" openapi:"$ref:SharedLimitation;type:array"`
	JobTitle                  string                                  `json:"jobTitle" openapi:"example:lorem ipsum"`
	CertificateCode           string                                  `json:"certificateCode" openapi:"example:1234567890"`
	Grace                     int                                     `json:"grace" openapi:"example:1"`
	Warning                   int                                     `json:"warning" openapi:"example:1"`
	Attention                 int                                     `json:"attention" openapi:"example:1"`
	Progress                  int                                     `json:"progress" openapi:"example:1"`
	ExperienceAmount          int                                     `json:"experienceAmount" openapi:"example:1"`
	ExperienceAmountUnit      string                                  `json:"experienceAmountUnit" openapi:"example:1"`
	IsSubcontractor           bool                                    `json:"isSubcontractor" openapi:"example:true"`
	CompanyRegistrationNumber string                                  `json:"companyRegistrationNumber" openapi:"example:1"`
	PercentLengthInContract   int                                     `json:"percentLengthInContract" openapi:"example:1"`
	HourLengthInContract      int                                     `json:"hourLengthInContract" openapi:"example:1"`
	Salary                    int                                     `json:"salary" openapi:"example:1"`
	StaffTypes                []StaffsQueryResponseDataStaffType      `json:"staffTypes" openapi:"$ref:StaffsQueryResponseDataStaffType;type:array"`
	TrialTime                 int                                     `json:"trial_time" openapi:"example:2022-01-01T00:00:00Z"`
	Attachments               []types.UploadMetadata                  `json:"attachments" openapi:"$ref:UploadMetadata;type:array;required"`
	JoinedAt                  int                                     `json:"joined_at" openapi:"example:2022-01-01T00:00:00Z"`
	ContractStartedAt         int                                     `json:"contract_started_at" openapi:"example:2022-01-01T00:00:00Z"`
	ContractExpiresAt         int                                     `json:"contract_expires_at" openapi:"example:2022-01-01T00:00:00Z"`
	CreatedAt                 int                                     `json:"created_at" openapi:"example:2022-01-01T00:00:00Z"`
	UpdatedAt                 int                                     `json:"updated_at" openapi:"example:2022-01-01T00:00:00Z"`
	DeletedAt                 int                                     `json:"deleted_at" openapi:"example:2022-01-01T00:00:00Z"`
}

/*
 * @apiDefine: StaffsQueryResponseData
 */
type StaffsQueryResponseData struct {
	Limit      int                           `json:"limit" openapi:"example:10"`
	Offset     int                           `json:"offset" openapi:"example:0"`
	Page       int                           `json:"page" openapi:"example:1"`
	TotalRows  int                           `json:"totalRows" openapi:"example:1"`
	TotalPages int                           `json:"totalPages" openapi:"example:1"`
	Items      []StaffsQueryResponseDataItem `json:"items" openapi:"$ref:StaffsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: StaffsQueryResponse
 */
type StaffsQueryResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       StaffsQueryResponseData `json:"data" openapi:"$ref:StaffsQueryResponseData"`
}

/*
 * @apiDefine: StaffsQueryNotFoundResponse
 */
type StaffsQueryNotFoundResponse struct {
	Staffs []domain.Staff `json:"staffs" openapi:"$ref:Staff;type:array"`
}
