package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCyclesDuplicateRequestBody_ValidateBody(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rBody := &CyclesDuplicateRequestBody{}
		body := strings.NewReader(`{"cycleId":1}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) > 0 {
			t.Errorf("Result should be successful but got errors: %#v\n", errs)
		}
	})
	t.Run("Failed validation", func(t *testing.T) {
		rBody := &CyclesDuplicateRequestBody{}
		body := strings.NewReader(`{"cycleId":0}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Errorf("Expect return error because body data is invalid errors: %#v\n", errs)
		}
	})
}
