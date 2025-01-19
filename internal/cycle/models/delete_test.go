package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCyclesDeleteRequestBody_ValidateBody(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rBody := &CyclesDeleteRequestBody{}
		body := strings.NewReader(`{"ids":[1,2,3]}`)
		r := httptest.NewRequest(http.MethodPost, "/", body)
		errs := rBody.ValidateBody(r)
		if len(errs) > 0 {
			t.Errorf("Result should be successful but got errors: %#v\n", errs)
		}
	})
}
