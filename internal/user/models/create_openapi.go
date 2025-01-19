package models

/*
 * @apiDefine: UsersCreateResponseDataLanguageSkill
 */
type UsersCreateResponseDataLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: UsersCreateResponseDataRolePermission
 */
type UsersCreateResponseDataRolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:test;required"`
	Title string `json:"title" openapi:"example:test;required"`
}

/*
 * @apiDefine: UsersCreateResponseDataRole
 */
type UsersCreateResponseDataRole struct {
	ID          uint                                    `json:"id" openapi:"example:1"`
	Name        string                                  `json:"name" openapi:"example:test;required"`
	Permissions []UsersCreateResponseDataRolePermission `json:"permissions" openapi:"$ref:UsersCreateResponseDataRolePermission;type:array;"`
}

/*
 * @apiDefine: UsersCreateResponseDataRoleID
 */
type UsersCreateResponseDataRoleID struct {
	ID uint `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: UsersCreateResponseData
 */
type UsersCreateResponseData struct {
	ID                      uint                                   `json:"id" openapi:"example:1"`
	RoleIDs                 []UsersCreateResponseDataRoleID        `json:"roleIds" openapi:"$ref:UsersCreateResponseDataRoleID;type:array"`
	StaffID                 uint                                   `json:"staffId" openapi:"example:1"`
	CustomerID              uint                                   `json:"customerId" openapi:"example:1"`
	FirstName               string                                 `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName                string                                 `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username                string                                 `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Password                string                                 `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Email                   string                                 `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone                   string                                 `json:"phone" openapi:"example:09123456789;required;"`
	Telephone               string                                 `json:"telephone" openapi:"example:02112345678"`
	LanguageSkills          []UsersCreateResponseDataLanguageSkill `json:"languageSkills" openapi:"$ref:UsersCreateResponseDataLanguageSkill;type:array;"`
	RegistrationNumber      string                                 `json:"registrationNumber" openapi:"example:1234567890"`
	WorkPhoneNumber         string                                 `json:"workPhoneNumber" openapi:"example:02112345678"`
	Gender                  string                                 `json:"gender" openapi:"example:male"`
	AccountNumber           string                                 `json:"accountNumber" openapi:"example:02112345678"`
	NationalCode            string                                 `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate               string                                 `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z"`
	AvatarUrl               string                                 `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	ForcedChangePassword    bool                                   `json:"forcedChangePassword" openapi:"example:true"`
	Roles                   []UsersCreateResponseDataRole          `json:"roles" openapi:"$ref:UsersCreateResponseDataRole;type:array"`
	PrivacyPolicyAcceptedAt string                                 `json:"privacy_policy_accepted_at" openapi:"example:2023-05-06T12:34:56Z"`
	SuspendedAt             string                                 `json:"suspended_at" openapi:"example:2023-05-06T12:34:56Z"`
	CreatedAt               string                                 `json:"created_at" openapi:"example:2023-05-06T12:34:56Z"`
	UpdatedAt               string                                 `json:"updated_at" openapi:"example:2023-05-06T12:34:56Z"`
	DeletedAt               string                                 `json:"deleted_at" openapi:"example:2023-05-06T12:34:56Z"`
}

/*
 * @apiDefine: UsersCreateResponse
 */
type UsersCreateResponse struct {
	StatusCode int                     `json:"statusCode" openapi:"example:200"`
	Data       UsersCreateResponseData `json:"data" openapi:"$ref:UsersCreateResponseData"`
}
