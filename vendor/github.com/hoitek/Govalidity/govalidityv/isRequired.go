package govalidityv

import (
	"github.com/hoitek/Govalidity/govaliditym"
	"github.com/hoitek/Govalidity/govalidityt"
	"strings"
)

func IsRequired(field string, params ...interface{}) (bool, error) {
	label, value := GetFieldLabelAndValue(field, params)
	err := GetErrorMessageByFieldValue(govaliditym.Default.IsRequired, label, value)
	switch convertedValue := value.(type) {
	case string:
		if strings.Trim(convertedValue, " ") == "" {
			return false, err
		}
	case *govalidityt.File:
		if convertedValue == nil {
			return false, err
		}
	case []*govalidityt.File:
		if len(convertedValue) == 0 {
			return false, err
		}
		for _, file := range convertedValue {
			if file == nil {
				return false, err
			}
		}
	}
	return true, nil
}
