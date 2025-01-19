package models

import (
	"testing"

	govalidity "github.com/hoitek/Govalidity"
)

func TestCyclesUpdateRequestParams_ValidateParams(t *testing.T) {
	t.Run("Success validation", func(t *testing.T) {
		rParams := &CyclesUpdateRequestParams{}
		params := govalidity.Params{
			"id": "1",
		}
		errs := rParams.ValidateParams(params)
		if len(errs) > 0 {
			t.Errorf("Expected no errors, got: %v", errs)
		}
	})
	// t.Run("Failed validation", func(t *testing.T) {
	// 	rParams := &CyclesUpdateRequestParams{}
	// 	params := govalidity.Params{}
	// 	errs := rParams.ValidateParams(params)
	// 	if len(errs) == 0 {
	// 		t.Errorf("Expected errors, got: %v", errs)
	// 	}
	// })
}
