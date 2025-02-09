package govalidityv

import (
	"github.com/hoitek/Govalidity/govaliditym"
	"net"
	"strings"
)

func IsIpV4(field string, params ...interface{}) (bool, error) {
	label, value := GetFieldLabelAndValue(field, params)
	err := GetErrorMessageByFieldValue(govaliditym.Default.IsIpV4, label, value)
	str := value.(string)
	ip := net.ParseIP(str)
	isValid := ip != nil && strings.Contains(str, ".")
	if isValid {
		return true, nil
	}
	return false, err
}
