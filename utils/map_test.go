package utils

import (
	"reflect"
	"testing"
)

type mapType = map[string]interface{}

func TestToMap(t *testing.T) {
	t.Run("When marshal has error", func(t *testing.T) {
		mapData := ToMap(mapType{
			"foo": make(chan int),
		})
		emptyMapData := mapType{}
		if !reflect.DeepEqual(mapData, emptyMapData) {
			t.Error("Expected empty map data")
		}
	})
}
