package utils_test

import (
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/utils"
	"testing"
)

// TestGetDBOperatorAndValue_ContainsOperatorWithLeadingAndTrailingSpaces tests the GetDBOperatorAndValue function with CONTAINS operator and leading/trailing spaces in the value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_ContainsOperatorWithLeadingAndTrailingSpaces(t *testing.T) {
	op := operators.CONTAINS
	value := " test "
	expected := utils.DBOperatorAndValue{
		Operator: "LIKE",
		Value:    "% test %",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_StartsWithOperatorWithLeadingAndTrailingSpaces tests the GetDBOperatorAndValue function with STARTS_WITH operator and leading/trailing spaces in the value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_StartsWithOperatorWithLeadingAndTrailingSpaces(t *testing.T) {
	op := operators.STARTS_WITH
	value := " test "
	expected := utils.DBOperatorAndValue{
		Operator: "LIKE",
		Value:    " test %",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_EndsWithOperatorWithLeadingAndTrailingSpaces tests the GetDBOperatorAndValue function with ENDS_WITH operator and leading/trailing spaces in the value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_EndsWithOperatorWithLeadingAndTrailingSpaces(t *testing.T) {
	op := operators.ENDS_WITH
	value := " test "
	expected := utils.DBOperatorAndValue{
		Operator: "LIKE",
		Value:    "% test ",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_IsNotEmptyOperator tests the GetDBOperatorAndValue function with IS_NOT_EMPTY operator.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_IsNotEmptyOperator(t *testing.T) {
	op := operators.IS_NOT_EMPTY
	value := "test"
	expected := utils.DBOperatorAndValue{
		Operator: "IS NOT NULL",
		Value:    "",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_IsAnyOfOperatorWithMultipleValues tests the GetDBOperatorAndValue function with IS_ANY_OF operator and multiple values.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_IsAnyOfOperatorWithMultipleValues(t *testing.T) {
	op := operators.IS_ANY_OF
	value := "test1,test2,test3"
	expected := utils.DBOperatorAndValue{
		Operator: "IN",
		Value:    "('test1','test2','test3')",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_IsAnyOfOperatorWithSingleValue tests the GetDBOperatorAndValue function with IS_ANY_OF operator and a single value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_IsAnyOfOperatorWithSingleValue(t *testing.T) {
	op := operators.IS_ANY_OF
	value := "test"
	expected := utils.DBOperatorAndValue{
		Operator: "IN",
		Value:    "('test')",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_DefaultOperatorForUnknownOperator tests the GetDBOperatorAndValue function with an unknown operator.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_DefaultOperatorForUnknownOperator(t *testing.T) {
	op := "unknown_operator"
	value := "test"
	expected := utils.DBOperatorAndValue{
		Operator: "=",
		Value:    value,
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_IsEmptyOperatorWithEmptyValue tests the GetDBOperatorAndValue function with IS_EMPTY operator and an empty value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_IsEmptyOperatorWithEmptyValue(t *testing.T) {
	op := operators.IS_EMPTY
	value := ""
	expected := utils.DBOperatorAndValue{
		Operator: "IS NULL",
		Value:    "",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}

// TestGetDBOperatorAndValue_IsNotEmptyOperatorWithNonEmptyValue tests the GetDBOperatorAndValue function with IS_NOT_EMPTY operator and a non-empty value.
//
// Parameter(s): t *testing.T
func TestGetDBOperatorAndValue_IsNotEmptyOperatorWithNonEmptyValue(t *testing.T) {
	op := operators.IS_NOT_EMPTY
	value := "test"
	expected := utils.DBOperatorAndValue{
		Operator: "IS NOT NULL",
		Value:    "",
	}
	result := utils.GetDBOperatorAndValue(op, value)
	if result.Operator != expected.Operator || result.Value != expected.Value {
		t.Errorf("Expected: %+v, got: %+v", expected, result)
	}
}
