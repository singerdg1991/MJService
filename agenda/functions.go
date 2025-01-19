package agenda

import (
	"strconv"
	"time"
)

// ConvertStringToInt64 converts string to int64
func ConvertStringToInt64(value string) (int64, error) {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// ConvertISODateTimeToTime converts ISO datetime to time.Time
func ConvertISODateTimeToTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// ConvertTimestampToISODateTime converts timestamp to ISO datetime
func ConvertTimestampToISODateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}

// GetCurrentTimeAsISODateTime returns current time as ISO datetime
func GetCurrentTimeAsISODateTime() string {
	return time.Now().Format(time.RFC3339)
}

// PrepareInterval prepares interval
func PrepareInterval(d time.Duration) string {
	return strconv.FormatInt(int64(d.Seconds()), 10)
}