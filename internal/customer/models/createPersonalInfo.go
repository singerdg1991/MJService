package models

import (
	"net/http"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
)

/*
 * @apiDefine: CustomersCreatePersonalInfoRequestBodyLimitation
 */
type CustomersCreatePersonalInfoRequestBodyLimitation struct {
	LimitationID uint   `json:"limitationId" openapi:"example:1"`
	Description  string `json:"description" openapi:"example:description"`
}

/*
 * @apiDefine: CustomersCreatePersonalInfoRequestBody
 */
type CustomersCreatePersonalInfoRequestBody struct {
	UserID               int64                                               `json:"-" openapi:"ignored"`
	DateOfBirthAsDate    *time.Time                                          `json:"-" openapi:"ignored"`
	StatusDateAsDate     *time.Time                                          `json:"-" openapi:"ignored"`
	MotherLangIDsInt64   []int64                                             `json:"-" openapi:"ignored"`
	FirstName            string                                              `json:"firstName" openapi:"example:John;required"`
	LastName             string                                              `json:"lastName" openapi:"example:Doe;required"`
	Gender               string                                              `json:"gender" openapi:"example:male"`
	DateOfBirth          string                                              `json:"dateOfBirth" openapi:"example:2021-01-01T00:00:00Z"`
	NationalCode         string                                              `json:"nationalCode" openapi:"example:1234567890"`
	MotherLangIDs        interface{}                                         `json:"motherLangIds" openapi:"example:[1,2,3];type:array;required;"`
	Email                string                                              `json:"email" openapi:"example:email@gmail.com"`
	PhoneNumber          string                                              `json:"phoneNumber" openapi:"example:123456789"`
	Status               string                                              `json:"status" openapi:"example:active"`
	StatusDate           string                                              `json:"statusDate" openapi:"example:2021-01-01T00:00:00Z"`
	Limitations          []*CustomersCreatePersonalInfoRequestBodyLimitation `json:"limitations" openapi:"$ref:CustomersCreatePersonalInfoRequestBodyLimitation;type:array"`
	Password             string                                              `json:"password" openapi:"example:123456;required;maxLen:100;minLen:2;"`
	ForcedChangePassword bool                                                `json:"forcedChangePassword" openapi:"example:true"`
	AuthenticatedUser    *sharedmodels.AuthenticatedUser                     `json:"-" openapi:"ignored"`
}

func (data *CustomersCreatePersonalInfoRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"firstName":            govalidity.New("firstName").MinMaxLength(2, 200).Required(),
		"lastName":             govalidity.New("lastName").MinMaxLength(2, 200).Required(),
		"gender":               govalidity.New("gender").In([]string{constants.GENDER_MALE, constants.GENDER_FEMALE}),
		"dateOfBirth":          govalidity.New("dateOfBirth").Required(),
		"nationalCode":         govalidity.New("nationalCode").Required(),
		"motherLangIds":        govalidity.New("motherLangIds"),
		"email":                govalidity.New("email").Email().Required(),
		"phoneNumber":          govalidity.New("phoneNumber").Number().Required(),
		"status":               govalidity.New("status").In([]string{constants.CUSTOMER_STATUS_ACTIVE, constants.CUSTOMER_STATUS_DEATH, constants.CUSTOMER_STATUS_FORMER_CUSTOMER_OR_DISCHARGED}).Required(),
		"statusDate":           govalidity.New("statusDate").Required(),
		"limitations":          govalidity.New("limitations").Optional(),
		"password":             govalidity.New("password").MinMaxLength(3, 25).Required(),
		"forcedChangePassword": govalidity.New("forcedChangePassword"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate dateOfBirth
	if data.DateOfBirth != "" {
		dateOfBirth, err := time.Parse(time.RFC3339, data.DateOfBirth)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"dateOfBirth": []string{"Date of birth is not valid"},
			}
		}
		data.DateOfBirthAsDate = &dateOfBirth
	}

	// Validate statusDate
	if data.StatusDate != "" {
		statusDate, err := time.Parse(time.RFC3339, data.StatusDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"statusDate": []string{"Status date is not valid"},
			}
		}
		data.StatusDateAsDate = &statusDate
	}

	return nil
}
