package ports

import "github.com/hoitek/Logger/engines"

type LoggerEngineType interface {
	*engines.LoggerEngineFile | *engines.LoggerEngineStdout | any
}
