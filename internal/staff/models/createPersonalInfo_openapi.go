package models

import "github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"

/*
 * @apiDefine: StaffsCreatePersonalInfoResponseDataLanguageSkill
 */
type StaffsCreatePersonalInfoResponseDataLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoResponseDataRolePermission
 */
type StaffsCreatePersonalInfoResponseDataRolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:test;required"`
	Title string `json:"title" openapi:"example:test;required"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoResponseDataRole
 */
type StaffsCreatePersonalInfoResponseDataRole struct {
	ID          uint                                                 `json:"id" openapi:"example:1"`
	Name        string                                               `json:"name" openapi:"example:test;required"`
	Permissions []StaffsCreatePersonalInfoResponseDataRolePermission `json:"permissions" openapi:"$ref:StaffsCreatePersonalInfoResponseDataRolePermission;type:array;"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoResponseDataRoleID
 */
type StaffsCreatePersonalInfoResponseDataRoleID struct {
	ID uint `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoResponseData
 */
type StaffsCreatePersonalInfoResponseData struct {
	ID                      uint                                                `json:"id" openapi:"example:1"`
	RoleIDs                 []StaffsCreatePersonalInfoResponseDataRoleID        `json:"roleIds" openapi:"$ref:StaffsCreatePersonalInfoResponseDataRoleID;type:array"`
	StaffID                 uint                                                `json:"staffId" openapi:"example:1"`
	CustomerID              uint                                                `json:"customerId" openapi:"example:1"`
	FirstName               string                                              `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName                string                                              `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username                string                                              `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Password                string                                              `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Email                   string                                              `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone                   string                                              `json:"phone" openapi:"example:09123456789;required;"`
	Telephone               string                                              `json:"telephone" openapi:"example:02112345678"`
	LanguageSkills          []StaffsCreatePersonalInfoResponseDataLanguageSkill `json:"languageSkills" openapi:"$ref:StaffsCreatePersonalInfoResponseDataLanguageSkill;type:array;"`
	RegistrationNumber      string                                              `json:"registrationNumber" openapi:"example:1234567890"`
	WorkPhoneNumber         string                                              `json:"workPhoneNumber" openapi:"example:02112345678"`
	Gender                  string                                              `json:"gender" openapi:"example:male"`
	AccountNumber           string                                              `json:"accountNumber" openapi:"example:02112345678"`
	NationalCode            string                                              `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate               string                                              `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z"`
	AvatarUrl               string                                              `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	ForcedChangePassword    bool                                                `json:"forcedChangePassword" openapi:"example:true"`
	Roles                   []StaffsCreatePersonalInfoResponseDataRole          `json:"roles" openapi:"$ref:StaffsCreatePersonalInfoResponseDataRole;type:array"`
	PrivacyPolicyAcceptedAt string                                              `json:"privacy_policy_accepted_at" openapi:"example:2023-05-06T12:34:56Z"`
	UserType                string                                              `json:"userType" openapi:"example:staff"`
	VehicleTypes            interface{}                                         `json:"vehicleTypes" openapi:"example:[\"car\",\"bicycle\",\"public_transportation\"];type:array;required;"`
	VehicleLicenseTypes     interface{}                                         `json:"vehicleLicenseTypes" openapi:"example:[\"automatic\",\"manual\"];type:array;required;"`
	Limitations             []sharedmodels.SharedLimitation                     `json:"limitations" openapi:"$ref:SharedLimitation;type:array;"`
	SuspendedAt             string                                              `json:"suspended_at" openapi:"example:2023-05-06T12:34:56Z"`
	CreatedAt               string                                              `json:"created_at" openapi:"example:2023-05-06T12:34:56Z"`
	UpdatedAt               string                                              `json:"updated_at" openapi:"example:2023-05-06T12:34:56Z"`
	DeletedAt               string                                              `json:"deleted_at" openapi:"example:2023-05-06T12:34:56Z"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoResponse
 */
type StaffsCreatePersonalInfoResponse struct {
	StatusCode int                                  `json:"statusCode" openapi:"example:200"`
	Data       StaffsCreatePersonalInfoResponseData `json:"data" openapi:"$ref:StaffsCreatePersonalInfoResponseData"`
}
