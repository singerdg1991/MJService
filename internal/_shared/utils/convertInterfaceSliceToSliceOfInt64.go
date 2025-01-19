package utils

import (
	"errors"
)

func ConvertInterfaceSliceToSliceOfInt64(numInterface interface{}) ([]int64, error) {
	idsInterface, ok := numInterface.([]interface{})
	if !ok {
		return nil, errors.New("input is invalid")
	}
	var ids []int64
	for _, id := range idsInterface {
		idFloat64, ok := id.(float64)
		if !ok {
			return nil, errors.New("input is invalid")
		}
		idInt64 := int64(idFloat64)
		if idInt64 <= 0 {
			return nil, errors.New("input is invalid")
		}
		ids = append(ids, idInt64)
	}
	return ids, nil
}

func ConvertInterfaceSliceToSliceOfString(input interface{}) ([]string, error) {
	idsInterface, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("input is invalid")
	}
	var ids []string
	for _, id := range idsInterface {
		idString, ok := id.(string)
		if !ok {
			return nil, errors.New("input is invalid")
		}
		if idString == "" {
			return nil, errors.New("input is invalid")
		}
		ids = append(ids, idString)
	}
	return ids, nil
}
