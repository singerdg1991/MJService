package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"reflect"
	"strconv"
)

func ToMap[T interface{}](data T) map[string]interface{} {
	bytes, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{}
	}
	dataMap := map[string]interface{}{}
	err = json.Unmarshal(bytes, &dataMap)
	if err != nil {
		return map[string]interface{}{}
	}
	return dataMap
}

func ToNumber(value interface{}) (res *float64, err error) {
	val := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		result := float64(val.Int())
		res = &result
	case uint, uint8, uint16, uint32, uint64:
		result := float64(val.Uint())
		res = &result
	case float32, float64:
		result := float64(val.Float())
		res = &result
	case string:
		result, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			res = nil
		} else {
			res = &result
		}
	default:
		err = fmt.Errorf("ToInt: unknown interface type %T", value)
		res = nil
	}

	return
}

func getJsonSlice(str string) []string {
	var strSlice []string
	isValid := json.Unmarshal([]byte(str), &strSlice) == nil
	if isValid {
		return strSlice
	}
	return []string{}
}

func ParseSQLOperator(op string, value interface{}) *filters.OperatorValue {
	sqlOp, ok := operators.SQL[op]
	if !ok {
		sqlOp = "="
	}

	if value == nil {
		return &filters.OperatorValue{
			Op:    sqlOp,
			Value: "",
		}
	}

	var strVal string = fmt.Sprintf("%v", value)
	var val interface{} = strVal

	switch op {
	case operators.CONTAINS, operators.STARTS_WITH, operators.ENDS_WITH:
		val = "%" + strVal + "%"
	case operators.IS_ANY_OF:
		val = getJsonSlice(strVal)
	}

	return &filters.OperatorValue{
		Op:    sqlOp,
		Value: val,
	}
}
