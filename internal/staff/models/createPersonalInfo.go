package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"net/http"
	"strings"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govaliditym"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	staffconstants "github.com/hoitek/Maja-Service/internal/staff/constants"
)

/*
 * @apiDefine: StaffsCreatePersonalInfoRequestBodyLanguageSkill
 */
type StaffsCreatePersonalInfoRequestBodyLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:test;required"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoRequestBodyLimitation
 */
type StaffsCreatePersonalInfoRequestBodyLimitation struct {
	LimitationID uint   `json:"limitationId" openapi:"example:1"`
	Description  string `json:"description" openapi:"example:description"`
}

/*
 * @apiDefine: StaffsCreatePersonalInfoRequestBody
 */
type StaffsCreatePersonalInfoRequestBody struct {
	LanguageSkillIds                 []int64                                          `json:"-" openapi:"ignored"`
	LanguageSkillsInt64              []int64                                          `json:"-" openapi:"ignored"`
	FirstName                        string                                           `json:"firstName" openapi:"example:saeed;required;maxLen:100;minLen:8;"`
	LastName                         string                                           `json:"lastName" openapi:"example:taher;required;maxLen:100;minLen:3;"`
	Username                         string                                           `json:"username" openapi:"example:saeed;required;maxLen:100;minLen:2;"`
	Password                         string                                           `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	Gender                           string                                           `json:"gender" openapi:"$ref:GenderEnum;type:string"`
	WorkPhoneNumber                  string                                           `json:"workPhoneNumber" openapi:"example:02112345678"`
	AccountNumber                    string                                           `json:"accountNumber" openapi:"example:02112345678"`
	LanguageSkills                   interface{}                                      `json:"languageSkills" openapi:"example:[1,2,3];type:array;required;"`
	ForcedChangePassword             bool                                             `json:"forcedChangePassword" openapi:"example:true"`
	Email                            string                                           `json:"email" openapi:"example:saeed@gmail.com;required;email;"`
	Phone                            string                                           `json:"phone" openapi:"example:09123456789;required;"`
	NationalCode                     string                                           `json:"nationalCode" openapi:"example:1234567890;required;"`
	BirthDate                        string                                           `json:"birthDate" openapi:"example:2023-05-06T12:34:56Z;required;date;"`
	AvatarUrl                        string                                           `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png;nullable;url;"`
	VehicleTypes                     interface{}                                      `json:"vehicleTypes" openapi:"example:[\"car\",\"bicycle\",\"public_transportation\"];type:array;required;"`
	VehicleLicenseTypes              interface{}                                      `json:"vehicleLicenseTypes" openapi:"example:[\"automatic\",\"manual\"];type:array;required;"`
	Limitations                      []*StaffsCreatePersonalInfoRequestBodyLimitation `json:"limitations" openapi:"$ref:StaffsCreatePersonalInfoRequestBodyLimitation;type:array"`
	VehicleTypesAsStringSlice        []string                                         `json:"-" openapi:"ignored"`
	VehicleLicenseTypesAsStringSlice []string                                         `json:"-" openapi:"ignored"`
}

func (data *StaffsCreatePersonalInfoRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"firstName":            govalidity.New("firstName").MinMaxLength(2, 200).Required(),
		"lastName":             govalidity.New("lastName").MinMaxLength(2, 200).Required(),
		"username":             govalidity.New("username").MinMaxLength(3, 25).Required(),
		"password":             govalidity.New("password").MinMaxLength(3, 25).Required(),
		"gender":               govalidity.New("gender").In([]string{constants.GENDER_MALE, constants.GENDER_FEMALE}),
		"email":                govalidity.New("email").Email().MinMaxLength(3, 25).Required(),
		"phone":                govalidity.New("phone").Required(),
		"workPhoneNumber":      govalidity.New("workPhoneNumber").Required(),
		"accountNumber":        govalidity.New("accountNumber").Optional(),
		"languageSkills":       govalidity.New("languageSkills"),
		"forcedChangePassword": govalidity.New("forcedChangePassword"),
		"nationalCode":         govalidity.New("nationalCode").Required(),
		"birthDate":            govalidity.New("birthDate").Required(),
		"avatarUrl":            govalidity.New("avatarUrl").Optional(),
		"vehicleTypes":         govalidity.New("vehicleTypes"),
		"vehicleLicenseTypes":  govalidity.New("vehicleLicenseTypes"),
		"limitations":          govalidity.New("limitations").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"firstName":            "First name",
			"lastName":             "Last name",
			"username":             "Username",
			"password":             "Password",
			"gender":               "Gender",
			"email":                "Email",
			"phone":                "Phone",
			"workPhoneNumber":      "Work phone number",
			"accountNumber":        "Account number",
			"languageSkills":       "Language skills",
			"forcedChangePassword": "Forced change password",
			"nationalCode":         "National Code",
			"birthDate":            "Birth date",
			"avatarUrl":            "Avatar url",
			"vehicleTypes":         "Vehicle types",
			"vehicleLicenseTypes":  "Vehicle license types",
			"limitations":          "Limitations",
		},
	)

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Sanitize phone and work phone number
	data.Phone = strings.ReplaceAll(data.Phone, " ", "")
	data.WorkPhoneNumber = strings.ReplaceAll(data.WorkPhoneNumber, " ", "")

	languageSkillsIDs, err := utils.ConvertInterfaceSliceToSliceOfInt64(data.LanguageSkills)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"languageSkills": []string{"Language skills must be an array of integers"},
		}
	}
	data.LanguageSkillsInt64 = languageSkillsIDs

	// Validate vehicle types
	vehicleTypes, err := utils.ConvertInterfaceSliceToSliceOfString(data.VehicleTypes)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"vehicleTypes": []string{"Vehicle types must be an array of strings"},
		}
	}
	data.VehicleTypesAsStringSlice = vehicleTypes
	for _, vehicleType := range vehicleTypes {
		if vehicleType != staffconstants.STAFF_VEHICLE_TYPE_CAR && vehicleType != staffconstants.STAFF_VEHICLE_TYPE_BICYCLE && vehicleType != staffconstants.STAFF_VEHICLE_TYPE_PUBLIC_TRANSPORTATION {
			return govalidity.ValidityResponseErrors{
				"vehicleTypes": []string{"Vehicle types must contain one of these values: " + staffconstants.STAFF_VEHICLE_TYPE_CAR + ", " + staffconstants.STAFF_VEHICLE_TYPE_BICYCLE + ", " + staffconstants.STAFF_VEHICLE_TYPE_PUBLIC_TRANSPORTATION},
			}
		}

	}

	// Validate vehicle license types
	vehicleLicenseTypes, err := utils.ConvertInterfaceSliceToSliceOfString(data.VehicleLicenseTypes)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"vehicleLicenseTypes": []string{"Vehicle license types must be an array of strings"},
		}
	}
	data.VehicleLicenseTypesAsStringSlice = vehicleLicenseTypes
	for _, vehicleLicenseType := range vehicleLicenseTypes {
		if vehicleLicenseType != staffconstants.STAFF_VEHICLE_LICENSE_TYPE_AUTOMATIC && vehicleLicenseType != staffconstants.STAFF_VEHICLE_LICENSE_TYPE_MANUAL {
			return govalidity.ValidityResponseErrors{
				"vehicleLicenseTypes": []string{"Vehicle license types must contain one of these values: " + staffconstants.STAFF_VEHICLE_LICENSE_TYPE_AUTOMATIC + ", " + staffconstants.STAFF_VEHICLE_LICENSE_TYPE_MANUAL},
			}
		}
	}

	return nil
}
