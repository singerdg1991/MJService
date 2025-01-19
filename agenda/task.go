package agenda

import (
	"fmt"
	"time"
)

// Task struct
type Task[T TaskValue] struct {
	Name string
	Function TaskFunction
	Type TaskType
	Value T
	LastRun time.Time
}

// NewTask creates new task
func NewTask[T TaskValue](name string, function TaskFunction, value T) *Task[T] {
	taskType := TYPE_INTERVAL
	if fmt.Sprintf("%T", value) == "*Task[string]" {
		taskType = TYPE_DATETIME
	}
	return &Task[T]{
		Name: name,
		Function: function,
		Type: taskType,
		Value: value,
	}
}