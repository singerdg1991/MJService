package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"net/http"
	"time"
)

/*
 * @apiDefine: CustomersUpdateAdditionalInfoRequestParams
 */
type CustomersUpdateAdditionalInfoRequestParams struct {
	UserID int `json:"userid,string" openapi:"example:1;nullable;pattern:^[0-9]+$;in:path"`
}

func (data *CustomersUpdateAdditionalInfoRequestParams) ValidateParams(params govalidity.Params) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"userid": govalidity.New("userid").Int().Required(),
	}

	errs := govalidity.ValidateParams(params, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

/*
 * @apiDefine: CustomersUpdateAdditionalInfoRequestBody
 */
type CustomersUpdateAdditionalInfoRequestBody struct {
	UserID                                             int64                           `json:"-" openapi:"ignored"`
	CustomerID                                         int64                           `json:"-" openapi:"ignored"`
	RelativeIDsInt64                                   []int64                         `json:"-" openapi:"ignored"`
	DiagnoseIDsInt64                                   []int64                         `json:"-" openapi:"ignored"`
	LimitingTheRightToSelfDeterminationStartDateAsDate *time.Time                      `json:"-" openapi:"ignored"`
	LimitingTheRightToSelfDeterminationEndDateAsDate   *time.Time                      `json:"-" openapi:"ignored"`
	SectionIDsInt64                                    []int64                         `json:"-" openapi:"ignored"`
	HasLimitingTheRightToSelfDetermination             bool                            `json:"hasLimitingTheRightToSelfDetermination" openapi:"ignored"`
	RelativeIDs                                        interface{}                     `json:"relativeIds" openapi:"example:[1,2,3]"`
	KeyNo                                              string                          `json:"keyNo" openapi:"example:1234567890"`
	PaymentMethod                                      string                          `json:"paymentMethod" openapi:"example:own"`
	SectionIDs                                         interface{}                     `json:"sectionIds" openapi:"example:[1,2,3];type:array;required;"`
	DiagnoseIDs                                        interface{}                     `json:"diagnoseIds" openapi:"example:[1,2,3]"`
	NurseGenderWish                                    string                          `json:"staffGenderWish" openapi:"example:male"`
	ResponsibleNurseID                                 *int                            `json:"responsibleStaffId" openapi:"example:null"`
	LimitingTheRightToSelfDeterminationStartDate       string                          `json:"limitingTheRightToSelfDeterminationStartDate" openapi:"example:2020-01-01"`
	LimitingTheRightToSelfDeterminationEndDate         string                          `json:"limitingTheRightToSelfDeterminationEndDate" openapi:"example:2020-01-01"`
	MobilityContract                                   string                          `json:"mobilityContract" openapi:"example:mobilityContract"`
	ParkingInfo                                        string                          `json:"parkingInfo" openapi:"example:parkingInfo"`
	ExtraExplanation                                   string                          `json:"extraExplanation" openapi:"example:extraExplanation"`
	AuthenticatedUser                                  *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *CustomersUpdateAdditionalInfoRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"relativeIds": govalidity.New("relativeIds").Optional(), //
		"keyNo":       govalidity.New("keyNo").Required(),
		"paymentMethod": govalidity.New("paymentMethod").In([]string{
			constants.PAYMENT_METHOD_OWN,
			constants.PAYMENT_METHOD_SETELI,
			constants.PAYMENT_METHOD_SETELI_AND_OWN,
		}),
		"sectionIds":         govalidity.New("sectionIds"),
		"diagnoseIds":        govalidity.New("diagnoseIds").Optional(), //
		"staffGenderWish":    govalidity.New("staffGenderWish").In([]string{constants.GENDER_MALE, constants.GENDER_FEMALE}).Required(),
		"responsibleStaffId": govalidity.New("responsibleStaffId").Int().Optional(),
		"limitingTheRightToSelfDeterminationStartDate": govalidity.New("limitingTheRightToSelfDeterminationStartDate"),
		"limitingTheRightToSelfDeterminationEndDate":   govalidity.New("limitingTheRightToSelfDeterminationEndDate"),
		"mobilityContract":                             govalidity.New("mobilityContract").Optional(),
		"parkingInfo":                                  govalidity.New("parkingInfo").Optional(),
		"extraExplanation":                             govalidity.New("extraExplanation").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate LimitingTheRightToSelfDeterminationStartDate
	if data.LimitingTheRightToSelfDeterminationStartDate != "" {
		limitingTheRightToSelfDeterminationStartDate, err := time.Parse("2006-01-02", data.LimitingTheRightToSelfDeterminationStartDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"limitingTheRightToSelfDeterminationStartDate": []string{"LimitingTheRightToSelfDeterminationStartDate is not valid"},
			}
		}
		data.LimitingTheRightToSelfDeterminationStartDateAsDate = &limitingTheRightToSelfDeterminationStartDate
	}

	// Validate LimitingTheRightToSelfDeterminationEndDate
	if data.LimitingTheRightToSelfDeterminationEndDate != "" {
		limitingTheRightToSelfDeterminationEndDate, err := time.Parse("2006-01-02", data.LimitingTheRightToSelfDeterminationEndDate)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"limitingTheRightToSelfDeterminationEndDate": []string{"LimitingTheRightToSelfDeterminationEndDate is not valid"},
			}
		}
		data.LimitingTheRightToSelfDeterminationEndDateAsDate = &limitingTheRightToSelfDeterminationEndDate
	}

	if data.LimitingTheRightToSelfDeterminationStartDateAsDate != nil || data.LimitingTheRightToSelfDeterminationEndDateAsDate != nil {
		data.HasLimitingTheRightToSelfDetermination = true
	}

	return nil
}
