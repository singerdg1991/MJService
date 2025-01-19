package utils

import (
	"fmt"
	"strings"

	"github.com/hoitek/Go-Quilder/operators"
)

type DBOperatorAndValue struct {
	Operator string
	Value    string
}

// GetDBOperatorAndValue returns the database operator and value for a given operator and value
func GetDBOperatorAndValue(op, value string) DBOperatorAndValue {
	switch op {
	case operators.CONTAINS:
		return DBOperatorAndValue{
			Operator: "LIKE",
			Value:    "%" + value + "%",
		}
	case operators.STARTS_WITH:
		return DBOperatorAndValue{
			Operator: "LIKE",
			Value:    value + "%",
		}
	case operators.ENDS_WITH:
		return DBOperatorAndValue{
			Operator: "LIKE",
			Value:    "%" + value,
		}
	case operators.EQUALS:
		return DBOperatorAndValue{
			Operator: "=",
			Value:    value,
		}
	case operators.IS_EMPTY:
		return DBOperatorAndValue{
			Operator: "IS NULL",
			Value:    "",
		}
	case operators.IS_NOT_EMPTY:
		return DBOperatorAndValue{
			Operator: "IS NOT NULL",
			Value:    "",
		}
	case operators.IS_ANY_OF:
		strSlice := strings.Split(value, ",")
		for i, v := range strSlice {
			strSlice[i] = "'" + v + "'"
		}
		return DBOperatorAndValue{
			Operator: "IN",
			Value:    fmt.Sprintf("(%s)", strings.Join(strSlice, ",")),
		}
	case operators.NUMBER_LESS_THAN_EQUALS:
		return DBOperatorAndValue{
			Operator: "<=",
			Value:    value,
		}
	}
	return DBOperatorAndValue{
		Operator: "=",
		Value:    value,
	}
}
