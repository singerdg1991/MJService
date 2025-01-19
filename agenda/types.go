package agenda

import "time"

// TaskType type of task type
type TaskType string

// TypeInterval type of interval
type TypeInterval time.Duration

// TaskFunction type
type TaskFunction func()
