package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/constants"
	"net/http"
	"time"
)

/*
 * @apiDefine: CustomersCreateServicesRequestBody
 */
type CustomersCreateServicesRequestBody struct {
	CustomerID           int        `json:"customerId" openapi:"example:1"`
	ServiceID            int        `json:"serviceId" openapi:"example:1"`
	ServiceTypeID        int        `json:"serviceTypeId" openapi:"example:1"`
	GradeID              int        `json:"gradeId" openapi:"example:1"`
	NurseWishID          int        `json:"staffWishId" openapi:"example:1"`
	ReportType           string     `json:"reportType" openapi:"example:reportType"`
	TimeValue            string     `json:"timeValue" openapi:"example:00:00:00"`
	Repeat               string     `json:"repeat" openapi:"example:weekly"`
	VisitType            string     `json:"visitType" openapi:"example:online"`
	ServiceLengthMinute  int        `json:"serviceLengthMinute" openapi:"example:60"`
	StartTimeRange       string     `json:"startTimeRange" openapi:"example:00:00:00"`
	EndTimeRange         *string    `json:"endTimeRange" openapi:"example:00:00:00"`
	Description          *string    `json:"description" openapi:"example:description"`
	PaymentMethod        string     `json:"paymentMethod" openapi:"example:own"`
	HomeCareFee          *int       `json:"homeCareFee" openapi:"example:100"`
	CityCouncilFee       *int       `json:"cityCouncilFee" openapi:"example:100"`
	TimeValueAsTime      *time.Time `json:"-" openapi:"ignored"`
	StartTimeRangeAsTime *time.Time `json:"-" openapi:"ignored"`
	EndTimeRangeAsTime   *time.Time `json:"-" openapi:"ignored"`
}

func (data *CustomersCreateServicesRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"customerId":    govalidity.New("customerId").Int().Required(),
		"serviceId":     govalidity.New("serviceId").Int().Required(),
		"serviceTypeId": govalidity.New("serviceTypeId").Int().Required(),
		"gradeId":       govalidity.New("gradeId").Int().Required(),
		"staffWishId":   govalidity.New("staffWishId").Int().Required(),
		"reportType":    govalidity.New("reportType").Required(),
		"timeValue":     govalidity.New("timeValue").Required(),
		"repeat": govalidity.New("repeat").In([]string{
			constants.SERVICE_REPEAT_DAILY,
			constants.SERVICE_REPEAT_WEEKLY,
			constants.SERVICE_REPEAT_MONTHLY,
			constants.SERVICE_REPEAT_EVERY_MONDAY,
			constants.SERVICE_REPEAT_EVERY_TUESDAY,
			constants.SERVICE_REPEAT_EVERY_WEDNESDAY,
			constants.SERVICE_REPEAT_EVERY_THURSDAY,
			constants.SERVICE_REPEAT_EVERY_FRIDAY,
			constants.SERVICE_REPEAT_EVERY_SATURDAY,
			constants.SERVICE_REPEAT_EVERY_SUNDAY,
		}).Required(),
		"visitType":           govalidity.New("visitType").In([]string{constants.VISIT_TYPE_ONLINE, constants.VISIT_TYPE_ONSITE}).Required(),
		"serviceLengthMinute": govalidity.New("serviceLengthMinute").Int().Min(1).Required(),
		"startTimeRange":      govalidity.New("startTimeRange").Required(),
		"endTimeRange":        govalidity.New("endTimeRange").Required(),
		"description":         govalidity.New("description"),
		"paymentMethod": govalidity.New("paymentMethod").In([]string{
			constants.PAYMENT_METHOD_OWN,
			constants.PAYMENT_METHOD_SETELI,
			constants.PAYMENT_METHOD_SETELI_AND_OWN,
		}),
		"homeCareFee":    govalidity.New("homeCareFee").Int().Min(0).Optional(),
		"cityCouncilFee": govalidity.New("cityCouncilFee").Int().Min(0).Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate timeValue
	timeValue, err := time.Parse("15:04:05", data.TimeValue)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"timeValue": []string{"Invalid time format"},
		}
	}
	data.TimeValueAsTime = &timeValue

	// Validate startTimeRange
	startTimeRange, err := time.Parse("15:04:05", data.StartTimeRange)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"startTimeRange": []string{"Invalid time format"},
		}
	}
	data.StartTimeRangeAsTime = &startTimeRange

	// Validate endTimeRange
	if data.EndTimeRange != nil {
		endTimeRange, err := time.Parse("15:04:05", *data.EndTimeRange)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"endTimeRange": []string{"Invalid time format"},
			}
		}
		data.EndTimeRangeAsTime = &endTimeRange
	}

	// Validate homeCareFee and cityCouncilFee
	if data.PaymentMethod == constants.PAYMENT_METHOD_OWN {
		if data.HomeCareFee == nil {
			return govalidity.ValidityResponseErrors{
				"homeCareFee": []string{"This field is required"},
			}
		}
		if data.CityCouncilFee != nil {
			return govalidity.ValidityResponseErrors{
				"cityCouncilFee": []string{"This field should be null"},
			}
		}
	}
	if data.PaymentMethod == constants.PAYMENT_METHOD_SETELI {
		if data.CityCouncilFee == nil {
			return govalidity.ValidityResponseErrors{
				"cityCouncilFee": []string{"This field is required"},
			}
		}
		if data.HomeCareFee != nil {
			return govalidity.ValidityResponseErrors{
				"homeCareFee": []string{"This field should be null"},
			}
		}
		data.HomeCareFee = data.CityCouncilFee
	}
	if data.PaymentMethod == constants.PAYMENT_METHOD_SETELI_AND_OWN {
		if data.CityCouncilFee == nil {
			return govalidity.ValidityResponseErrors{
				"cityCouncilFee": []string{"This field is required"},
			}
		}
		if data.HomeCareFee == nil {
			return govalidity.ValidityResponseErrors{
				"homeCareFee": []string{"This field is required"},
			}
		}
	}

	return nil
}
