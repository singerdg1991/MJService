package domain

import (
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	"testing"
	"time"
)

// Declare variables
var (
	tableName = "cycles"
	newCycle  = NewCycle()
)

// TestNewCycle tests the NewCycle function
func TestNewCycle(t *testing.T) {
	createdCycle := NewCycle()
	if createdCycle == nil {
		t.Error("NewCycle function is not working properly")
	}
}

// TestTableName tests the TableName function
func TestTableName(t *testing.T) {
	if tableName != newCycle.TableName() {
		t.Error("New cycle table name is not equal to created cycle table name")
	}
}

// TestToJson tests the ToJson method
func TestToJson(t *testing.T) {
	jsonString, err := newCycle.ToJson()
	if err != nil {
		t.Error("Error in to json method")
	}
	if jsonString == "" {
		t.Error("Json string is empty")
	}
}

// TestSetDefaultStatus tests the SetDefaultStatus method
func TestSetDefaultStatus(t *testing.T) {
	nextDay := time.Now().AddDate(0, 0, 1)
	newCycle.FreezePeriodDate = nextDay
	newCycle.EndDate = &nextDay
	newCycle.SetDefaultStatus()
	if newCycle.Status != constants.STATUS_ACTIVE {
		t.Error("Default status is not active")
	}

	// Check if cycle is frozen
	newCycle.FreezePeriodDate = time.Now().AddDate(0, 0, -1)
	if time.Now().After(newCycle.FreezePeriodDate) {
		newCycle.SetDefaultStatus()
		if newCycle.Status != constants.STATUS_FROZEN {
			t.Error("Default status is not frozen")
		}
	}

	// Check if cycle is expired
	prevDay := time.Now().AddDate(0, 0, -1)
	newCycle.FreezePeriodDate = nextDay
	newCycle.EndDate = &prevDay
	if newCycle.EndDate != nil && time.Now().After(*newCycle.EndDate) {
		newCycle.SetDefaultStatus()
		if newCycle.Status != constants.STATUS_EXPIRED {
			t.Error("Default status is not expired")
		}
	}
}
