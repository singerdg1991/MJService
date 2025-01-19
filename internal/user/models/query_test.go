package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersQueryRequestValidateQueries(t *testing.T) {
	t.Run("Succeed validation", func(t *testing.T) {
		rQuery := &UsersQueryRequestParams{}
		r := httptest.NewRequest(http.MethodPost, "/?id=1", nil)
		errs := rQuery.ValidateQueries(r)
		if len(errs) > 0 {
			t.Error("Can not return errors because query data is valid")
		}
	})

	t.Run("Failed validation when query is invalid", func(t *testing.T) {
		rQuery := &UsersQueryRequestParams{}
		r := httptest.NewRequest(http.MethodPost, "/?id=value", nil)
		errs := rQuery.ValidateQueries(r)
		if len(errs) == 0 {
			t.Error("Expect return error because query data is invalid")
		}
	})
}
