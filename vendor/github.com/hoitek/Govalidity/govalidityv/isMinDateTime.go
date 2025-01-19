package govalidityv

import (
	"github.com/hoitek/Govalidity/govaliditym"
	"time"
)

func IsMinDateTime(field string, value interface{}, min time.Time) (bool, error) {
	if value.(time.Time).Before(min) {
		return false, GetErrorMessageByFieldValue(govaliditym.Default.IsMinDateTime, field, value)
	}
	return true, nil
}
