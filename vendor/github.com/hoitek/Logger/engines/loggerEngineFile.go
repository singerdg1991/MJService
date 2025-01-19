package engines

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"time"
)

type LoggerEngineFile struct {
	Path string
}

func (l *LoggerEngineFile) Trace(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/trace_" + today + ".log"
	err := writeLog(filepath, args, "Trace")
	if err != nil {
		log.Println(err)
	}
}

func (l *LoggerEngineFile) Debug(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/debug" + today + ".log"
	err := writeLog(filepath, args, "Debug")
	if err != nil {
		log.Println(err)
	}
}

func (l *LoggerEngineFile) Info(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/info_" + today + ".log"
	err := writeLog(filepath, args, "Info")
	if err != nil {
		log.Println(err)
	}
}

func (l *LoggerEngineFile) Warn(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/warn_" + today + ".log"
	err := writeLog(filepath, args, "Warn")
	if err != nil {
		log.Println(err)
	}
}

func (l *LoggerEngineFile) Error(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/error_" + today + ".log"
	err := writeLog(filepath, args, "Error")
	if err != nil {
		log.Println(err)
	}
}

func (l *LoggerEngineFile) Fatal(args ...interface{}) {
	today := time.Now().Format("2006-01-02")
	filepath := l.Path + "/fatal_" + today + ".log"
	err := writeLog(filepath, args, "Fatal")
	if err != nil {
		log.Println(err)
	}
}

func writeLog(filepath string, data []interface{}, label string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var result map[string]interface{} = map[string]interface{}{}
	result["type"] = label
	result["time"] = time.Now()
	result["stack"] = string(debug.Stack())

	result["log"] = data
	buffer, err := json.Marshal(result)
	if err != nil {
		return err
	}
	buffer = append(buffer, 10)
	_, err = file.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
