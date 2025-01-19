package logger

import "github.com/hoitek/Logger/ports"

type LoggerEngine[T ports.LoggerEngineType] struct {
	Types     []string
	Instances *[]T
}
