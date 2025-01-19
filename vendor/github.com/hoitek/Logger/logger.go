package logger

import (
	"reflect"

	"github.com/hoitek/Logger/engines"
	"github.com/hoitek/Logger/ports"
)

type Logger[T ports.LoggerEngineType] struct {
	Engine LoggerEngine[T]
}

var Enable = true
var Default *Logger[ports.LoggerEngineType] = Initialize[ports.LoggerEngineType]()

func Initialize[T ports.LoggerEngineType](engines ...T) *Logger[T] {
	var loggerTypes []string = []string{}
	for _, engine := range engines {
		loggerType := reflect.ValueOf(&engine).Elem().Type().Name()
		loggerTypes = append(loggerTypes, loggerType)
	}
	var Default = &Logger[T]{}
	Default.Engine.Instances = &engines
	Default.Engine.Types = loggerTypes
	return Default
}
func (l Logger[T]) GetInstance() *Logger[ports.LoggerEngineType] {
	return Default
}

// Log with injected driver
func (l Logger[T]) Trace(args ...interface{}) {
	log("Trace", args...)
}
func (l Logger[T]) Debug(args ...interface{}) {
	log("Debug", args...)
}
func (l Logger[T]) Info(args ...interface{}) {
	log("Info", args...)
}
func (l Logger[T]) Warn(args ...interface{}) {
	log("Warn", args...)
}
func (l Logger[T]) Error(args ...interface{}) {
	log("Error", args...)
}
func (l Logger[T]) Fatal(args ...interface{}) {
	log("Fatal", args...)
}

// Log by Default driver
func Trace(args ...interface{}) {
	Default.Trace(args...)
}
func Debug(args ...interface{}) {
	Default.Debug(args...)
}
func Info(args ...interface{}) {
	Default.Info(args...)
}
func Warn(args ...interface{}) {
	Default.Warn(args...)
}
func Error(args ...interface{}) {
	Default.Error(args...)
}
func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

func log(functionName string, args ...interface{}) {
	if !Enable {
		return
	}
	for _, instance := range *Default.Engine.Instances {
		loggerFile, ok := instance.(*engines.LoggerEngineFile)
		if ok {
			method := getFunction(loggerFile, functionName)
			go method(args...)
		}
		loggerStdout, ok := instance.(*engines.LoggerEngineStdout)
		if ok {
			method := getFunction(loggerStdout, functionName)
			go method(args...)
		}
	}
}

func getFunction[T ports.LoggerEngineType](obj T, functionName string) func(...interface{}) {
	methodVal := reflect.ValueOf(obj).MethodByName(functionName)
	methodInterface := methodVal.Interface()
	method := methodInterface.(func(...interface{}))
	return method
}
