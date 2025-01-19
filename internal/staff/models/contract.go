package models

import (
	"encoding/json"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/config"
	"log"
	"net/http"
)

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestParams
 */
type StaffsCreateOrUpdateContractRequestParams struct {
	UserID int `json:"userId,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *StaffsCreateOrUpdateContractRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userId": govalidity.New("userId").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodySection
 */
type StaffsCreateOrUpdateContractRequestBodySection struct {
	ID int64 `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyAbility
 */
type StaffsCreateOrUpdateContractRequestBodyAbility struct {
	ID int64 `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyPaymentType
 */
type StaffsCreateOrUpdateContractRequestBodyPaymentType struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyContractType
 */
type StaffsCreateOrUpdateContractRequestBodyContractType struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyAvailableShifts
 */
type StaffsCreateOrUpdateContractRequestBodyAvailableShifts struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyStaffType
 */
type StaffsCreateOrUpdateContractRequestBodyStaffType struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBodyRole
 */
type StaffsCreateOrUpdateContractRequestBodyRole struct {
	ID int `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: StaffsCreateOrUpdateContractRequestBody
 */
type StaffsCreateOrUpdateContractRequestBody struct {
	ID                          int                                                      `json:"id" openapi:"ignored"`
	SectionIDs                  []int64                                                  `json:"sectionIds" openapi:"ignored"`
	UserID                      int                                                      `json:"userId" openapi:"ignored"`
	ShiftTypeIDs                []int64                                                  `json:"shiftTypeIDs" openapi:"ignored"`
	ContractTypeIDs             []int64                                                  `json:"contractTypeIds" openapi:"ignored"`
	StaffTypeIDs                []int64                                                  `json:"staffTypeIds" openapi:"ignored"`
	RoleIDs                     []int64                                                  `json:"roleIds" openapi:"ignored"`
	JobTitle                    string                                                   `json:"jobTitle" openapi:"example:jobTitle"`
	CertificateCode             string                                                   `json:"certificateCode" openapi:"example:1234567890"`
	JoinedAt                    string                                                   `json:"joinedAt" openapi:"type:string;example:2021-01-01T00:00:00Z"`
	ContractStartedAt           string                                                   `json:"contractStartedAt" openapi:"type:string;example:2021-01-01T00:00:00Z"`
	ContractExpiresAt           string                                                   `json:"contractExpiresAt" openapi:"type:string;example:2021-01-01T00:00:00Z"`
	TrialTime                   string                                                   `json:"trialTime" openapi:"type:string;example:2021-01-01T00:00:00Z"`
	OrganizationNumber          string                                                   `json:"organizationNumber" openapi:"ignored"`
	PercentLengthInContract     int                                                      `json:"percentLengthInContract,string" openapi:"type:string;example:100"`
	HourLengthInContract        int                                                      `json:"hourLengthInContract,string" openapi:"type:string;example:40"`
	Salary                      int                                                      `json:"salary,string" openapi:"type:string;example:1000"`
	ExperienceAmount            int                                                      `json:"experienceAmount,string" openapi:"type:string;example:1"`
	ExperienceAmountUnit        string                                                   `json:"experienceAmountUnit" openapi:"type:string;example:year"`
	CompanyRegistrationNumber   string                                                   `json:"companyRegistrationNumber" openapi:"type:string;example:1234567-8"`
	ShiftTypes                  *string                                                  `json:"availableShifts" openapi:"example:[{\"id\":1}]"`
	ContractTypes               *string                                                  `json:"contractTypes" openapi:"example:[{\"id\":1}]"`
	Sections                    *string                                                  `json:"sections" openapi:"example:[{\"id\":1}]"`
	StaffTypes                  *string                                                  `json:"staffTypes" openapi:"example:[{\"id\":1}]"`
	PaymentType                 interface{}                                              `json:"paymentType" openapi:"example:{\"id\":1}"`
	Roles                       *string                                                  `json:"roles" openapi:"example:[{\"id\":1}]"`
	Attachments                 []*govalidityt.File                                      `json:"attachments" openapi:"format:binary;type:array"`
	PreviousAttachments         string                                                   `json:"previousAttachments" openapi:"example:[{\"fileName\": \"424e5ebcf1e4b4f11707315705332860929.png\", \"fileSize\": 44547, \"path\": \"/uploads/staff\"}]"`
	PreviousAttachmentsMetadata []types.UploadMetadata                                   `json:"-" openapi:"ignored"`
	ShiftTypesAsMetadata        []StaffsCreateOrUpdateContractRequestBodyAvailableShifts `json:"-" openapi:"ignored"`
	ContractTypesAsMetadata     []StaffsCreateOrUpdateContractRequestBodyContractType    `json:"-" openapi:"ignored"`
	SectionsAsMetadata          []StaffsCreateOrUpdateContractRequestBodySection         `json:"-" openapi:"ignored"`
	StaffTypesAsMetadata        []StaffsCreateOrUpdateContractRequestBodyStaffType       `json:"-" openapi:"ignored"`
	PaymentTypeAsMetadata       StaffsCreateOrUpdateContractRequestBodyPaymentType       `json:"-" openapi:"ignored"`
	RolesAsMetadata             []StaffsCreateOrUpdateContractRequestBodyRole            `json:"-" openapi:"ignored"`
}

func (data *StaffsCreateOrUpdateContractRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"joinedAt":                  govalidity.New("joinedAt").Optional(),
		"contractStartedAt":         govalidity.New("contractStartedAt").Optional(),
		"contractExpiresAt":         govalidity.New("contractExpiresAt").Optional(),
		"availableShifts":           govalidity.New("availableShifts").Optional(),
		"trialTime":                 govalidity.New("trialTime").Optional(),
		"listOfContract":            govalidity.New("listOfContract").Optional(),
		"percentLengthInContract":   govalidity.New("percentLengthInContract").Optional(),
		"hourLengthInContract":      govalidity.New("hourLengthInContract").Optional(),
		"salary":                    govalidity.New("salary").Optional(),
		"sections":                  govalidity.New("sections").Optional(),
		"teams":                     govalidity.New("teams").Optional(),
		"abilities":                 govalidity.New("abilities").Optional(),
		"staffTypes":                govalidity.New("staffTypes").Optional(),
		"experienceAmount":          govalidity.New("experienceAmount").Optional(),
		"experienceAmountUnit":      govalidity.New("experienceAmountUnit").Optional(),
		"companyRegistrationNumber": govalidity.New("companyRegistrationNumber").Optional(),
		"paymentType":               govalidity.New("paymentType").Optional(),
		"jobTitle":                  govalidity.New("jobTitle").Optional(),
		"certificateCode":           govalidity.New("certificateCode").Optional(),
		"roles":                     govalidity.New("roles").Optional(),
		"contractTypes":             govalidity.New("contractTypes").Optional(),
		"attachments":               govalidity.New("attachments").Files(),
		"previousAttachments":       govalidity.New("previousAttachments").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate uploaded files size
	fileErrs := utils.ValidateUploadFilesSize("attachments", data.Attachments, config.StaffConfig.MaxUploadSizeLimit)
	if fileErrs != nil {
		return fileErrs
	}

	// Validate uploaded files mime type
	fileErrs = utils.ValidateUploadFilesMimeType("attachments", data.Attachments, []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
	})
	if fileErrs != nil {
		return fileErrs
	}

	// Try to unmarshal PreviousAttachmentsMetadata
	if data.PreviousAttachments != "" {
		if err := json.Unmarshal([]byte(data.PreviousAttachments), &data.PreviousAttachmentsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"previousAttachments": []string{"previousAttachments is not a valid JSON"},
			}
		}
	}

	// Validate ShiftTypes
	if data.ShiftTypes != nil {
		if err := json.Unmarshal([]byte(*data.ShiftTypes), &data.ShiftTypesAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"shiftTypes": []string{"shiftTypes is not a valid JSON"},
			}
		}
	}

	// Validate ContractTypes
	if data.ContractTypes != nil {
		if err := json.Unmarshal([]byte(*data.ContractTypes), &data.ContractTypesAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"contractTypes": []string{"contractTypes is not a valid JSON"},
			}
		}
	}

	// Validate Sections
	if data.Sections != nil {
		if err := json.Unmarshal([]byte(*data.Sections), &data.SectionsAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"sections": []string{"sections is not a valid JSON"},
			}
		}
	}

	// Validate StaffTypes
	if data.StaffTypes != nil {
		if err := json.Unmarshal([]byte(*data.StaffTypes), &data.StaffTypesAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"staffTypes": []string{"staffTypes is not a valid JSON"},
			}
		}
	}

	// Validate PaymentType
	if data.PaymentType != nil {
		log.Printf("PaymentType: %#v\n", data.PaymentType)
		convertedPaymentType := map[string]int{}
		paymentType, ok := data.PaymentType.(map[string]interface{})
		if !ok {
			return govalidity.ValidityResponseErrors{
				"paymentType": []string{"paymentType1 is not a valid JSON"},
			}
		}
		for k, v := range paymentType {
			convertedPaymentType[k] = int(v.(float64))
		}
		paymentTypeAsJson, err := json.Marshal(convertedPaymentType)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"paymentType": []string{"paymentType2 is not a valid JSON"},
			}
		}
		if err := json.Unmarshal(paymentTypeAsJson, &data.PaymentTypeAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"paymentType": []string{"paymentType3 is not a valid JSON"},
			}
		}
	}

	// Validate Roles
	if data.Roles != nil {
		if err := json.Unmarshal([]byte(*data.Roles), &data.RolesAsMetadata); err != nil {
			return govalidity.ValidityResponseErrors{
				"roles": []string{"roles is not a valid JSON"},
			}
		}
	}

	return nil
}
