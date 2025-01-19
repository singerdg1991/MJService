package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCyclesCreateRequestBody_ValidateBody(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rBody := &CyclesCreateRequestBody{}
		body := strings.NewReader(`{
			"sectionId":1,
			"start_date":"2025-01-01",
			"end_date":null,
			"periodLength":"oneWeek",
			"shiftMorningStartTime":"08:00",
			"shiftMorningEndTime":"12:00",
			"shiftEveningStartTime": "12:00",
			"shiftEveningEndTime": "18:00",
			"shiftNightStartTime": "18:00",
			"shiftNightEndTime": "00:00",
            "freeze_period_date": "2025-01-02",
			"wishDays": 4
		}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) > 0 {
			t.Errorf("Expect to be successful but errors: %#v\n", errs)
		}
	})
	t.Run("Failed validation when start date is invalid", func(t *testing.T) {
		rBody := &CyclesCreateRequestBody{}
		body := strings.NewReader(`{
			"sectionId":1,
            "start_date":"2025",
            "end_date":null,
            "periodLength":"oneWeek",
            "shiftMorningStartTime":"08:00",
            "shiftMorningEndTime":"12:00",
            "shiftEveningStartTime": "12:00",
            "shiftEveningEndTime": "18:00",
            "shiftNightStartTime": "18:00",
            "shiftNightEndTime": "00:00",
            "freeze_period_date": "2025-01-02",
            "wishDays": 4
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because start date is invalid")
		}
	})
	t.Run("Failed validation when start date is in the past", func(t *testing.T) {
		rBody := &CyclesCreateRequestBody{}
		body := strings.NewReader(`{
			"sectionId":1,
            "start_date":"2021-01-01",
            "end_date":null,
            "periodLength":"oneWeek",
            "shiftMorningStartTime":"08:00",
            "shiftMorningEndTime":"12:00",
            "shiftEveningStartTime": "12:00",
            "shiftEveningEndTime": "18:00",
            "shiftNightStartTime": "18:00",
            "shiftNightEndTime": "00:00",
            "freeze_period_date": "2025-01-02",
            "wishDays": 4
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because start date is in the past")
		}
	})
	t.Run("Failed validation when end date is prepare and has wrong format", func(t *testing.T) {
		rBody := &CyclesCreateRequestBody{}
		body := strings.NewReader(`{
			"sectionId":1,
            "start_date":"2025-01-01",
            "end_date":"2025",
            "periodLength":null,
            "shiftMorningStartTime":"08:00",
            "shiftMorningEndTime":"12:00",
            "shiftEveningStartTime": "12:00",
            "shiftEveningEndTime": "18:00",
            "shiftNightStartTime": "18:00",
            "shiftNightEndTime": "00:00",
            "freeze_period_date": "2025-01-02",
            "wishDays": 4
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because end date is prepare and has wrong format")
		}
	})
	t.Run("Failed validation when end date is prepare and is before start date", func(t *testing.T) {
		rBody := &CyclesCreateRequestBody{}
		body := strings.NewReader(`{
			"sectionId":1,
            "start_date":"2025-01-01",
            "end_date":"2024-01-01",
            "periodLength":null,
            "shiftMorningStartTime":"08:00",
            "shiftMorningEndTime":"12:00",
            "shiftEveningStartTime": "12:00",
            "shiftEveningEndTime": "18:00",
            "shiftNightStartTime": "18:00",
            "shiftNightEndTime": "00:00",
            "freeze_period_date": "2025-01-02",
            "wishDays": 4
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because end date is prepare and is before start date")
		}
	})
}
