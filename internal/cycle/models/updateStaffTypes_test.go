package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCyclesUpdateStaffTypesRequestBody_ValidateBody(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"2025-01-01",
                    "neededStaffCount":1,
                    "roleId":1,
                    "startHour":"00:00",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) > 0 {
			t.Errorf("Expected no errors, got: %v", errs)
		}
	})
	t.Run("When shiftName is wrong", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"wrong",
                    "datetime":"2025-01-01",
                    "neededStaffCount":1,
                    "roleId":1,
                    "startHour":"00:00",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
	t.Run("When neededStaffCount is less than zero", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"2025-01-01",
                    "neededStaffCount":-1,
                    "roleId":1,
                    "startHour":"00:00",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
	t.Run("When roleId is less than one", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"2025-01-01",
                    "neededStaffCount":1,
                    "roleId":0,
                    "startHour":"00:00",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
	t.Run("When datetime is wrong", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"wrong",
                    "neededStaffCount":1,
                    "roleId":1,
                    "startHour":"00:00",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
	t.Run("When startHour is wrong", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"2025-01-01",
                    "neededStaffCount":1,
                    "roleId":1,
                    "startHour":"wrong",
                    "endHour":"00:00"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
	t.Run("When endHour is wrong", func(t *testing.T) {
		rBody := &CyclesUpdateStaffTypesRequestBody{}
		body := strings.NewReader(`{
            "staffTypes": [
                {
                    "shiftName":"morning",
                    "datetime":"2025-01-01",
                    "neededStaffCount":1,
                    "roleId":1,
                    "startHour":"00:00",
                    "endHour":"wrong"
                }
            ]
        }`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expected errors, got: %v", errs)
		}
	})
}
