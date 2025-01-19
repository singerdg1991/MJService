package utils

import (
	"bytes"
	"encoding/csv"
)

// ParseCsvFile function to parse a csv file
func ParseCsvFile(content []byte) (map[string]interface{}, error) {
	r := csv.NewReader(bytes.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	for _, record := range records {
		if len(record) > 1 {
			data[record[0]] = record[1]
		}
	}
	return data, nil
}
