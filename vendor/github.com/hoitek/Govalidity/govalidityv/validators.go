package govalidityv

import (
	"errors"
	"fmt"
	"github.com/hoitek/Govalidity/govaliditym"
	"strings"
)

func GetFieldLabelAndValue(field string, params []interface{}) (string, interface{}) {
	fieldLabel := field
	value := params[0]
	label, ok := (*govaliditym.FieldLabels)[field]
	if ok {
		fieldLabel = label
	}
	return fieldLabel, value
}

func GetErrorMessageByFieldValue(baseErrorMessage string, field string, value interface{}) error {
	valueType := fmt.Sprintf("%T", value)
	errMessage := strings.ReplaceAll(baseErrorMessage, "{field}", field)
	if valueType != "*govalidityt.File" && valueType != "[]*govalidityt.File" {
		errMessage = strings.ReplaceAll(errMessage, "{value}", value.(string))
	}
	return errors.New(errMessage)
}
