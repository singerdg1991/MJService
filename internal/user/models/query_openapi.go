package models

import "github.com/hoitek/Maja-Service/internal/user/domain"

/*
 * @apiDefine: UsersQueryResponseDataItemLanguageSkill
 */
type UsersQueryResponseDataItemLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: UsersQueryResponseDataItemRolePermission
 */
type UsersQueryResponseDataItemRolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:test;required"`
	Title string `json:"title" openapi:"example:test;required"`
}

/*
 * @apiDefine: UsersQueryResponseDataItemRole
 */
type UsersQueryResponseDataItemRole struct {
	ID          uint                                       `json:"id" openapi:"example:1"`
	Name        string                                     `json:"name" openapi:"example:test;required"`
	Permissions []UsersQueryResponseDataItemRolePermission `json:"permissions" openapi:"$ref:UsersQueryResponseDataItemRolePermission;type:array;"`
}

/*
 * @apiDefine: UsersQueryResponseDataItemRoleID
 */
type UsersQueryResponseDataItemRoleID struct {
	ID uint `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: UsersQueryResponseDataItem
 */
type UsersQueryResponseDataItem struct {
	ID                      uint                                      `json:"id" openapi:"example:1"`
	RoleIDs                 []UsersQueryResponseDataItemRoleID        `json:"roleIds" openapi:"$ref:UsersQueryResponseDataItemRoleID;type:array"`
	StaffID                 uint                                      `json:"staffId" openapi:"example:1"`
	CustomerID              uint                                      `json:"customerId" openapi:"example:1"`
	FirstName               string                                    `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName                string                                    `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username                string                                    `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Password                string                                    `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Email                   string                                    `json:"email" openapi:"example:sgh370@yahoo.com"`
	Phone                   string                                    `json:"phone" openapi:"example:09123456789;required;"`
	Telephone               string                                    `json:"telephone" openapi:"example:02112345678"`
	LanguageSkills          []UsersQueryResponseDataItemLanguageSkill `json:"languageSkills" openapi:"$ref:UsersQueryResponseDataItemLanguageSkill;type:array;"`
	RegistrationNumber      string                                    `json:"registrationNumber" openapi:"example:1234567890"`
	WorkPhoneNumber         string                                    `json:"workPhoneNumber" openapi:"example:02112345678"`
	Gender                  string                                    `json:"gender" openapi:"example:male"`
	AccountNumber           string                                    `json:"accountNumber" openapi:"example:02112345678"`
	NationalCode            string                                    `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate               string                                    `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z"`
	AvatarUrl               string                                    `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	ForcedChangePassword    bool                                      `json:"forcedChangePassword" openapi:"example:true"`
	Roles                   []UsersQueryResponseDataItemRole          `json:"roles" openapi:"$ref:UsersQueryResponseDataItemRole;type:array"`
	PrivacyPolicyAcceptedAt string                                    `json:"privacy_policy_accepted_at" openapi:"example:2023-05-06T12:34:56Z"`
	SuspendedAt             string                                    `json:"suspended_at" openapi:"example:2023-05-06T12:34:56Z"`
	CreatedAt               string                                    `json:"created_at" openapi:"example:2023-05-06T12:34:56Z"`
	UpdatedAt               string                                    `json:"updated_at" openapi:"example:2023-05-06T12:34:56Z"`
	DeletedAt               string                                    `json:"deleted_at" openapi:"example:2023-05-06T12:34:56Z"`
}

/*
 * @apiDefine: UsersQueryResponseData
 */
type UsersQueryResponseData struct {
	Limit      int                          `json:"limit" openapi:"example:10"`
	Offset     int                          `json:"offset" openapi:"example:0"`
	Page       int                          `json:"page" openapi:"example:1"`
	TotalRows  int                          `json:"totalRows" openapi:"example:1"`
	TotalPages int                          `json:"totalPages" openapi:"example:1"`
	Items      []UsersQueryResponseDataItem `json:"items" openapi:"$ref:UsersQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: UsersQueryResponse
 */
type UsersQueryResponse struct {
	StatusCode int                    `json:"statusCode" openapi:"example:200"`
	Data       UsersQueryResponseData `json:"data" openapi:"$ref:UsersQueryResponseData"`
}

/*
 * @apiDefine: UsersQueryNotFoundResponse
 */
type UsersQueryNotFoundResponse struct {
	Users []domain.User `json:"users" openapi:"$ref:User;type:array"`
}
