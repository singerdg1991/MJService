package utils_test

import (
	"github.com/hoitek/Maja-Service/utils"
	"testing"
	"time"
)

// TestTryParseToDateTime tests the TryParseToDateTime function.
//
// It defines test cases with input, expected output, and error.
// Verifies that the actual output matches the expected output.
func TestTryParseToDateTime(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		input    interface{}
		expected time.Time
		err      error
	}{
		{"2022-01-01T00:00:00Z", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), nil},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call the TryParseToDateTime function
		actual, err := utils.TryParseToDateTime(tc.input)

		// Check that the actual and expected values are equal
		if actual != tc.expected || err != tc.err {
			t.Errorf("TryParseToDateTime(%v) = (%v, %v), want (%v, %v)", tc.input, actual, err, tc.expected, tc.err)
		}
	}
}
