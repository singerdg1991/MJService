package engines

import (
	"log"
)

type LoggerEngineStdout struct {
}

func (l *LoggerEngineStdout) Trace(args ...interface{}) {
	var result []interface{}
	result = append(result, "Trace: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}

func (l *LoggerEngineStdout) Debug(args ...interface{}) {
	var result []interface{}
	result = append(result, "Debug: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}

func (l *LoggerEngineStdout) Info(args ...interface{}) {
	var result []interface{}
	result = append(result, "Info: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}

func (l *LoggerEngineStdout) Warn(args ...interface{}) {
	var result []interface{}
	result = append(result, "Warning: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}

func (l *LoggerEngineStdout) Error(args ...interface{}) {
	var result []interface{}
	result = append(result, "Error: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}

func (l *LoggerEngineStdout) Fatal(args ...interface{}) {
	var result []interface{}
	result = append(result, "Fatal: ")
	for _, arg := range args {
		result = append(result, arg)
	}
	log.Println(result...)
}
