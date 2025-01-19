package agenda

import "time"

// TaskValue type of value of task
type TaskValue interface {
	time.Duration | string | any
}