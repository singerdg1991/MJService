package domain_test

import (
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"testing"
)

// TestCycleChat_TableName tests the TableName function of the CycleChat struct.
//
// t is a pointer to testing.T, which is used to manage test state and support logging of failures.
// A boolean indicating whether the test was successful.
func TestCycleChat_TableName(t *testing.T) {
	expected := "cycleChats"
	result := (&domain.CycleChat{}).TableName()
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}
}
