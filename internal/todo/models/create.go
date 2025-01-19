package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/todo/domain"
	"net/http"
	"time"
)

/*
 * @apiDefine: TodosCreateRequestBody
 */
type TodosCreateRequestBody struct {
	Title             string           `json:"title" openapi:"example:title;required;maxLen:100;minLen:2;"`
	Date              string           `json:"date" openapi:"example:2021-01-01;required;"`
	Time              string           `json:"time" openapi:"example:00:00;required;"`
	UserID            uint             `json:"userId" openapi:"example:1;required;"`
	Description       string           `json:"description" openapi:"example:description;required;maxLen:100;minLen:2;"`
	DateAsDate        *time.Time       `json:"-" openapi:"ignored"`
	TimeAsTime        *time.Time       `json:"-" openapi:"ignored"`
	User              *domain.TodoUser `json:"-" openapi:"ignored"`
	AuthenticatedUser *domain.TodoUser `json:"-" openapi:"ignored"`
}

func (data *TodosCreateRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":       govalidity.New("title").MinMaxLength(3, 255).Required(),
		"date":        govalidity.New("date").MinMaxLength(3, 255).Required(),
		"time":        govalidity.New("time").MinMaxLength(3, 255).Required(),
		"userId":      govalidity.New("userId").Int().Min(1).Required(),
		"description": govalidity.New("description").MinMaxLength(3, 2000).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Check date
	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"date": []string{"Invalid date format, should be YYYY-MM-DD"},
		}
	}
	data.DateAsDate = &date

	// Check time should be in 24 hour format
	time, err := time.Parse("15:04", data.Time)
	if err != nil {
		return govalidity.ValidityResponseErrors{
			"time": []string{"Invalid time format, should be HH:MM"},
		}
	}
	data.TimeAsTime = &time

	return nil
}
