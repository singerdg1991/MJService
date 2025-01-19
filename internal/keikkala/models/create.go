package models

import (
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/keikkala/constants"
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
)

/*
 * @apiDefine: KeikkalasCreateRequestBody
 */
type KeikkalasCreateRequestBody struct {
	StartDate         string                          `json:"start_date" openapi:"example:2021-01-01;required"`
	EndDate           string                          `json:"end_date" openapi:"example:2021-01-01;required"`
	StartTime         string                          `json:"start_time" openapi:"example:08:00;required"`
	EndTime           string                          `json:"end_time" openapi:"example:16:00;required"`
	RoleID            int                             `json:"roleId" openapi:"example:1;required"`
	KaupunkiAddress   *string                         `json:"kaupunkiAddress" openapi:"example:address;required"`
	SectionIDs        interface{}                     `json:"sectionIds" openapi:"example:[1,2,3];type:array;required;"`
	PaymentType       string                          `json:"paymentType" openapi:"example:paySoon;required"`
	Description       *string                         `json:"description" openapi:"example:description;required"`
	StartDateAsDate   *time.Time                      `json:"-" openapi:"ignored"`
	EndDateAsDate     *time.Time                      `json:"-" openapi:"ignored"`
	StartTimeAsTime   *time.Time                      `json:"-" openapi:"ignored"`
	EndTimeAsTime     *time.Time                      `json:"-" openapi:"ignored"`
	StartDateAndTime  *time.Time                      `json:"-" openapi:"ignored"`
	EndDateAndTime    *time.Time                      `json:"-" openapi:"ignored"`
	SectionIDsInt64   []int64                         `json:"-" openapi:"ignored"`
	AuthenticatedUser *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *KeikkalasCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"start_date":      govalidity.New("start_date").Required(),
		"end_date":        govalidity.New("end_date").Required(),
		"start_time":      govalidity.New("start_time").Required(),
		"end_time":        govalidity.New("end_time").Required(),
		"roleId":          govalidity.New("roleId").Int().Min(1).Required(),
		"kaupunkiAddress": govalidity.New("kaupunkiAddress").Optional(),
		"sectionIds":      govalidity.New("sectionIds"),
		"paymentType": govalidity.New("paymentType").In([]string{
			constants.KEIKKALA_PAYMENT_TYPE_PAYSOON,
			constants.KEIKKALA_PAYMENT_TYPE_BONUS,
			constants.KEIKKALA_PAYMENT_TYPE_NOTHING,
		}).Required(),
		"description": govalidity.New("description").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate startDate
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"start_date must be a valid date in format YYYY-MM-DD"},
		}
	}
	data.StartDateAsDate = &startDate

	// Validate endDate
	endDate, err := time.Parse("2006-01-02", data.EndDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"end_date": []string{"end_date must be a valid date in format YYYY-MM-DD"},
		}
	}
	data.EndDateAsDate = &endDate

	// Validate startTime
	startTime, err := time.Parse("15:04", data.StartTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_time": []string{"start_time must be a valid time in format HH:MM"},
		}
	}
	data.StartTimeAsTime = &startTime

	// Validate endTime
	endTime, err := time.Parse("15:04", data.EndTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"end_time": []string{"end_time must be a valid time in format HH:MM"},
		}
	}
	data.EndTimeAsTime = &endTime

	// Combine startDate and startTime
	startDateAndTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), startTime.Location())
	data.StartDateAndTime = &startDateAndTime

	// Combine endDate and endTime
	endDateAndTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), endTime.Location())
	data.EndDateAndTime = &endDateAndTime

	// Convert interface slice to slice of int64
	sectionIDs, err := utils.ConvertInterfaceSliceToSliceOfInt64(data.SectionIDs)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"sectionIds": []string{"sectionIds must be an array of integers"},
		}
	}
	data.SectionIDsInt64 = sectionIDs

	return nil
}
