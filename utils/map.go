package utils

import "encoding/json"

func ToMap[T interface{}](data T) map[string]interface{} {
	bytes, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{}
	}
	dataMap := map[string]interface{}{}
	err = json.Unmarshal(bytes, &dataMap)
	if err != nil {
		return map[string]interface{}{}
	}
	return dataMap
}
