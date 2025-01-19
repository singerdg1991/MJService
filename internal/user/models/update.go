package models

import (
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govaliditym"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	roleDomain "github.com/hoitek/Maja-Service/internal/role/domain"
)

/*
 * @apiDefine: UsersUpdateRequestParams
 */
type UsersUpdateRequestParams struct {
	ID int `json:"id,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *UsersUpdateRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id": govalidity.New("ID").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)

	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: UsersUpdateRequestBodyLanguageSkill
 */
type UsersUpdateRequestBodyLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: UsersUpdateRequestBody
 */
type UsersUpdateRequestBody struct {
	UserType             string                                `json:"userType" openapi:"$ref:UserTypeEnum;type:string;required;example:staff"`
	LanguageSkillIds     []int64                               `json:"-" openapi:"ignored"`
	FirstName            string                                `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName             string                                `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username             string                                `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Password             string                                `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Gender               string                                `json:"gender" openapi:"$ref:GenderEnum;type:string"`
	WorkPhoneNumber      string                                `json:"workPhoneNumber" openapi:"example:02112345678"`
	AccountNumber        string                                `json:"accountNumber" openapi:"example:02112345678"`
	LanguageSkills       []UsersUpdateRequestBodyLanguageSkill `json:"languageSkills" openapi:"$ref:UsersUpdateRequestBodyLanguageSkill;type:array;"`
	RegistrationNumber   string                                `json:"registrationNumber" openapi:"example:1234567890"`
	ForcedChangePassword bool                                  `json:"forcedChangePassword" openapi:"example:true"`
	Email                string                                `json:"email" openapi:"example:saeed@gmail.com;required;email;"`
	Phone                string                                `json:"phone" openapi:"example:09123456789;required;"`
	NationalCode         string                                `json:"nationalCode" openapi:"example:1234567890;required;"`
	BirthDate            string                                `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z;required;date;"`
	AvatarUrl            string                                `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	Role                 *roleDomain.Role                      `json:"role" openapi:"ignored"`
}

func (data *UsersUpdateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"firstName":            govalidity.New("firstName").MinMaxLength(2, 200).Required(),
		"lastName":             govalidity.New("lastName").MinMaxLength(2, 200).Required(),
		"username":             govalidity.New("username").MinMaxLength(3, 25).Required(),
		"password":             govalidity.New("password"),
		"gender":               govalidity.New("gender").In([]string{constants.GENDER_MALE, constants.GENDER_FEMALE}),
		"userType":             govalidity.New("userType").In([]string{constants.USER_TYPE_CUSTOMER, constants.USER_TYPE_STAFF}).Required(),
		"email":                govalidity.New("email").Email().MinMaxLength(3, 25).Required(),
		"phone":                govalidity.New("phone").Required(),
		"workPhoneNumber":      govalidity.New("workPhoneNumber").Required(),
		"accountNumber":        govalidity.New("accountNumber").Required(),
		"languageSkills":       govalidity.New("languageSkills").Required(),
		"registrationNumber":   govalidity.New("registrationNumber").Required(),
		"forcedChangePassword": govalidity.New("forcedChangePassword"),
		"nationalCode":         govalidity.New("nationalCode").Required(),
		"birthDate":            govalidity.New("birthDate").Required(),
		"avatarUrl":            govalidity.New("avatarUrl").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"firstName":            "First name",
			"lastName":             "Last name",
			"username":             "Username",
			"password":             "Password",
			"gender":               "Gender",
			"userType":             "User Type",
			"email":                "Email",
			"phone":                "Phone",
			"workPhoneNumber":      "Work phone number",
			"accountNumber":        "Account number",
			"languageSkills":       "Language skills",
			"registrationNumber":   "Registration number",
			"forcedChangePassword": "Forced change password",
			"nationalCode":         "National code",
			"birthDate":            "Birth date",
			"avatarUrl":            "Avatar url",
		},
	)

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
