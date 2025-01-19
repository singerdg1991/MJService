package utils

import (
	"encoding/json"
	"strings"
	"sync"
)

// GetMapValueFromDotJoinedKeys function to get value from a map[string]interface{} by a dot joined key
func GetMapValueFromDotJoinedKeys(data map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	value := data
	for _, subKey := range keys {
		val, ok := value[subKey].(map[string]interface{})
		if !ok {
			return value[subKey]
		} else {
			value = val
		}
	}
	return ""
}

// JoinMapKeysWithDot function to join map keys with dot
func JoinMapKeysWithDot(m map[string]interface{}, prefix string, keys *[]string) {
	for k, v := range m {
		newPrefix := prefix + k
		switch t := v.(type) {
		case map[string]interface{}:
			JoinMapKeysWithDot(t, newPrefix+".", keys)
		default:
			*keys = append(*keys, newPrefix)
		}
	}
}

// ToJson function to convert a map[string]interface{} to json string
func ToJson(data map[string]interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SetMapValueBasedOnNestedKeys(data map[string]interface{}, keys []string, value interface{}, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(keys) == 1 {
		mx.Lock()
		data[keys[0]] = value
		mx.Unlock()
		return
	}
	mx.Lock()
	if _, ok := data[keys[0]]; !ok {
		data[keys[0]] = make(map[string]interface{})
		mx.Unlock()
		return
	}
	nestedMap := data[keys[0]].(map[string]interface{})
	mx.Unlock()
	wg.Add(1)
	go SetMapValueBasedOnNestedKeys(nestedMap, keys[1:], value, mx, wg)
}
