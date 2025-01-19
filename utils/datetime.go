package utils

import (
	"fmt"
	"time"
)

// TryParseToDateTime attempts to parse a given interface to a time.Time.
func TryParseToDateTime(t interface{}) (time.Time, error) {
	// Convert the interface to a string
	str := fmt.Sprintf("%v", t)

	// Define the layout for parsing the string to time.Time
	layout := time.RFC3339

	// Attempt to parse the string to time.Time
	return time.Parse(layout, str)
}
