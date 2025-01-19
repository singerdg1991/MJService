package govalidityv

import (
	"github.com/hoitek/Govalidity/govalidityconv"
	"github.com/hoitek/Govalidity/govaliditym"
)

func IsNumber(field string, params ...interface{}) (bool, error) {
	label, value := GetFieldLabelAndValue(field, params)
	err := GetErrorMessageByFieldValue(govaliditym.Default.IsNumber, label, value)
	number, numberError := govalidityconv.ToNumber(value)
	if numberError != nil || number == nil {
		return false, err
	}
	return true, nil
}
