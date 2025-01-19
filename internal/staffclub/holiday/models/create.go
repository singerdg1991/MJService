package models

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/constants"
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: HolidaysCreateRequestBody
 */
type HolidaysCreateRequestBody struct {
	StartDate         string                          `json:"start_date" openapi:"example:2021-01-01;required"`
	EndDate           string                          `json:"end_date" openapi:"example:2021-01-01;required"`
	Title             string                          `json:"title" openapi:"example:title;required"`
	PaymentType       *string                         `json:"paymentType" openapi:"example:withSalary;required"`
	Description       *string                         `json:"description" openapi:"example:description;required"`
	UserID            *int                            `json:"userId" openapi:"example:1;required"`
	StartDateAsDate   time.Time                       `json:"-" openapi:"ignored"`
	EndDateAsDate     time.Time                       `json:"-" openapi:"ignored"`
	AuthenticatedUser sharedmodels.AuthenticatedUser  `json:"-" openapi:"ignored"`
	User              *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *HolidaysCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"start_date":  govalidity.New("start_date").Required(),
		"end_date":    govalidity.New("end_date").Required(),
		"title":       govalidity.New("title").MinMaxLength(3, 255).Required(),
		"paymentType": govalidity.New("paymentType"),
		"description": govalidity.New("description").MinMaxLength(3, 255).Optional(),
		"userId":      govalidity.New("userId").Int().Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Convert start_date and end_date to time.Time
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"start_date": []string{"start_date must be a valid date"},
		}
	}
	data.StartDateAsDate = startDate

	endDate, err := time.Parse("2006-01-02", data.EndDate)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"end_date": []string{"end_date must be a valid date"},
		}
	}
	data.EndDateAsDate = endDate

	// Check paymentType
	if data.PaymentType != nil {
		if *data.PaymentType != constants.HOLIDAY_PAYMENT_TYPE_WITH_SALARY && *data.PaymentType != constants.HOLIDAY_PAYMENT_TYPE_WITHOUT_SALARY {
			return govalidity.ValidityResponseErrors{
				"paymentType": []string{fmt.Sprintf("paymentType must be one of %s or %s", constants.HOLIDAY_PAYMENT_TYPE_WITH_SALARY, constants.HOLIDAY_PAYMENT_TYPE_WITHOUT_SALARY)},
			}
		}
	} else {
		data.PaymentType = new(string)
		*data.PaymentType = constants.HOLIDAY_PAYMENT_TYPE_WITH_SALARY
	}

	return nil
}
