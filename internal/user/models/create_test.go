package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUsersCreateRequestBodyValidateBody(t *testing.T) {
	t.Run("Failed validation", func(t *testing.T) {
		rBody := &UsersCreateRequestBody{}
		body := strings.NewReader(`{"namo":"value"}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because body data is invalid")
		}
	})

	t.Run("Failed validation", func(t *testing.T) {
		rBody := &UsersCreateRequestBody{}
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		errs := rBody.ValidateBody(r)
		if len(errs) == 0 {
			t.Error("Expect return error because body data is invalid")
		}
	})
}
