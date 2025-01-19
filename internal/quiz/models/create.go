package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/quiz/constants"
	"net/http"
	"time"
)

/*
 * @apiDefine: QuizzesCreateRequestBody
 */
type QuizzesCreateRequestBody struct {
	Title                     string      `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	StartDateTime             string      `json:"startDateTime" openapi:"example:2021-01-01T00:00:00Z;required;"`
	EndDateTime               *string     `json:"endDateTime" openapi:"example:2021-01-01T00:00:00Z;required;"`
	DurationInMinute          *int        `json:"durationInMinute" openapi:"example:60;required;"`
	Status                    string      `json:"status" openapi:"example:disable;required;"`
	AvailableParticipantType  *string     `json:"availableParticipantType" openapi:"example:all;"`
	ParticipantUserIDs        interface{} `json:"participantUserIDs" openapi:"example:[1,2,3];"`
	ParticipantUserIDsAsInt64 []int64     `json:"-" openapi:"ignored"`
	IsLock                    string      `json:"isLock" openapi:"example:false;required;"`
	IsLockAsBool              *bool       `json:"-" openapi:"ignored"`
	Password                  *string     `json:"password" openapi:"example:password;required;"`
	Description               *string     `json:"description" openapi:"example:description;required;"`
	StartDateTimeAsDate       *time.Time  `json:"-" openapi:"ignored"`
	EndDateTimeAsDate         *time.Time  `json:"-" openapi:"ignored"`
}

func (data *QuizzesCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":                    govalidity.New("title").MinMaxLength(3, 25).Required(),
		"startDateTime":            govalidity.New("startDateTime").Required(),
		"endDateTime":              govalidity.New("endDateTime"),
		"durationInMinute":         govalidity.New("durationInMinute"),
		"status":                   govalidity.New("status").In([]string{constants.QUIZ_STATUS_ENABLE, constants.QUIZ_STATUS_DISABLE}).Required(),
		"availableParticipantType": govalidity.New("availableParticipantType").In([]string{constants.QUIZ_AVAILABLE_PARTICIPANT_TYPE_ALL, constants.QUIZ_AVAILABLE_PARTICIPANT_TYPE_CUSTOMER, constants.QUIZ_AVAILABLE_PARTICIPANT_TYPE_STAFF}).Required(),
		"participantUserIDs":       govalidity.New("participantUserIDs"),
		"isLock":                   govalidity.New("isLock").In([]string{"true", "false"}).Required(),
		"password":                 govalidity.New("password"),
		"description":              govalidity.New("description"),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Validate startDateTime
	startDateTime, err := time.Parse(time.RFC3339, data.StartDateTime)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"startDateTime": []string{"Start datetime is invalid"},
		}
	}
	data.StartDateTimeAsDate = &startDateTime

	// Validate endDateTime if not nil
	if data.EndDateTime != nil {
		endDateTime, err := time.Parse(time.RFC3339, *data.EndDateTime)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"endDateTime": []string{"End datetime is invalid"},
			}
		}
		data.EndDateTimeAsDate = &endDateTime
	}

	// Validate durationInMinute if not nil
	if data.DurationInMinute != nil {
		if *data.DurationInMinute <= 0 {
			return govalidity.ValidityResponseErrors{
				"durationInMinute": []string{"Duration in minute must be greater than 0"},
			}
		}
	}

	// Validate isLock
	isLockAsBool := false
	if data.IsLock == "true" {
		isLockAsBool = true
	}
	data.IsLockAsBool = &isLockAsBool

	// Validate password with isLock is true
	if *data.IsLockAsBool && data.Password == nil {
		return govalidity.ValidityResponseErrors{
			"password": []string{"Password is required"},
		}
	}

	return nil
}
