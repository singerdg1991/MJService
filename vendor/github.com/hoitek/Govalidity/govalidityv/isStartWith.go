package govalidityv

import (
	"errors"
	"fmt"
	"github.com/hoitek/Govalidity/govalidityconv"
	"github.com/hoitek/Govalidity/govaliditym"
	"strings"
)

func IsStartWith(field string, dataValue interface{}, str string) (bool, error) {
	fieldLabel := field

	label, ok := (*govaliditym.FieldLabels)[field]
	if ok {
		fieldLabel = label
	}

	value := ""
	number, errConv := govalidityconv.ToNumber(dataValue)
	if errConv == nil && number != nil {
		value = fmt.Sprintf("%v", *number)
	} else {
		value = dataValue.(string)
	}

	if strings.Index(value, str) == 0 {
		return true, nil
	}

	errMessage := strings.ReplaceAll(govaliditym.Default.IsMaxLength, "{field}", fieldLabel)
	errMessage = strings.ReplaceAll(errMessage, "{value}", value)
	errMessage = strings.ReplaceAll(errMessage, "{str}", str)
	err := errors.New(errMessage)

	return false, err
}
