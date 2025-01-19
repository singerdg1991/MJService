package utils

import "encoding/json"

func Jsonify(v interface{}) (interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return "", err
	}
	return f, nil
}
