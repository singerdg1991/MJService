package models

import (
	"net/http"
	"time"

	govalidity "github.com/hoitek/Govalidity"
	cdomain "github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
)

/*
 * @apiDefine: CyclesCreateVisitUnplannedRequestBody
 */
type CyclesCreateVisitUnplannedRequestBody struct {
	CycleID       int                           `json:"cycleId" openapi:"example:1;"`
	Date          string                        `json:"date" openapi:"example:2021-01-01;required;"`
	Time          string                        `json:"time" openapi:"example:00:00;required;"`
	Length        int                           `json:"length" openapi:"example:1;required;"`
	CustomerID    int                           `json:"customerId" openapi:"example:1;required;"`
	TaskType      string                        `json:"taskType" openapi:"example:taskType;required;"`
	ServiceID     int                           `json:"serviceId" openapi:"example:1;required;"`
	ServiceTypeID int                           `json:"serviceTypeId" openapi:"example:1;required;"`
	StaffID       int                           `json:"staffId" openapi:"example:1;required;"`
	Description   *string                       `json:"description" openapi:"example:This is a description;"`
	Customer      *cdomain.Customer             `json:"-" openapi:"ignored"`
	Staff         *domain.CyclePickupShiftStaff `json:"-" openapi:"ignored"`
	TimeAsTime    *time.Time                    `json:"-" openapi:"ignored"`
	DateAsDate    *time.Time                    `json:"-" openapi:"ignored"`
}

// ValidateBody validates the body of an HTTP request against a predefined schema.
//
// The function takes an HTTP request as a parameter and checks its body against a schema defined for the CyclesCreateVisitUnplannedRequestBody struct.
// It returns a govalidity.ValidityResponseErrors object containing any validation errors that occurred.
func (data *CyclesCreateVisitUnplannedRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"cycleId":       govalidity.New("cycleId").Int().Min(1).Required(),
		"date":          govalidity.New("date").Required(),
		"time":          govalidity.New("time").Required(),
		"length":        govalidity.New("length").Int().Min(1).Required(),
		"customerId":    govalidity.New("customerId").Int().Min(1).Required(),
		"taskType":      govalidity.New("taskType").Required(),
		"serviceId":     govalidity.New("serviceId").Int().Min(1).Required(),
		"serviceTypeId": govalidity.New("serviceTypeId").Int().Min(1).Required(),
		"staffId":       govalidity.New("staffId").Int().Min(1).Required(),
		"description":   govalidity.New("description").Optional(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate time
	t, err := time.Parse("15:04", data.Time)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"time": []string{"time must be in the format of HH:MM"},
		}
	}
	data.TimeAsTime = &t

	// Validate date
	d, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"date": []string{"date must be in the format of YYYY-MM-DD"},
		}
	}
	data.DateAsDate = &d

	return nil
}
