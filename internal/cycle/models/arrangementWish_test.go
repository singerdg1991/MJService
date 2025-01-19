package models

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCyclesArrangementWishRequestBody_ValidateBody(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rBody := &CyclesArrangementWishRequestBody{}
		body := strings.NewReader(`{"cycleId":1}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) != 0 {
			log.Printf("%#v\n", errs)
			t.Error("Error is not nil")
		}
	})
	t.Run("Failed validation when cycleId is 0", func(t *testing.T) {
		rBody := &CyclesArrangementWishRequestBody{}
		body := strings.NewReader(`{"cycleId":0}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because body data is invalid")
		}
	})
	t.Run("Failed validation when cycleId is minus", func(t *testing.T) {
		rBody := &CyclesArrangementWishRequestBody{}
		body := strings.NewReader(`{"cycleId":-1}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because body data is invalid")
		}
	})
}
