package govalidityv

import (
	"github.com/hoitek/Govalidity/govaliditym"
	"strings"
)

func IsLowerCase(field string, params ...interface{}) (bool, error) {
	label, value := GetFieldLabelAndValue(field, params)
	err := GetErrorMessageByFieldValue(govaliditym.Default.IsLowerCase, label, value)
	str := value.(string)
	isValid := str == strings.ToLower(str)
	if isValid {
		return true, nil
	}
	return false, err
}
