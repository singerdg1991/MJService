package agenda

import (
	"log"
	"time"
)

// Manager struct
type Manager struct {
	Tasks []any

}

// Setup new manager
func Setup() *Manager {
	return &Manager{}
}

// AddTask new task
func (m *Manager) AddTask(task any) *Manager {
	switch assertedTask := task.(type) {
	case *Task[time.Duration]:
		assertedTask.Type = TYPE_INTERVAL
	case *Task[string]:
		assertedTask.Type = TYPE_DATETIME
	}
	m.Tasks = append(m.Tasks, task)
	return m
}

// Run all tasks
func (m *Manager) Run() {
	go func() {
		for {
			for _, task := range m.Tasks {
				go func(task any) {
					switch assertedTask := task.(type) {
					case *Task[time.Duration]:
						// Convert task.Value to int64
						intervalValue, err := ConvertStringToInt64(PrepareInterval(assertedTask.Value))
						if err != nil {
							log.Printf("Task %s has invalid interval value %s with the error: %s\n", assertedTask.Name, assertedTask.Value, err.Error())
						}
						// Check if current time is greater than last run time + interval value
						if time.Now().Unix() > assertedTask.LastRun.Unix() + intervalValue {
							// Run task function
							log.Printf("Task %s is running\n", assertedTask.Name)
							go assertedTask.Function()
							// Update last run time
							assertedTask.LastRun = time.Now()
						}
					case *Task[string]:
						// Convert ISO Datetime to time.Time
						datetimeValue, err := ConvertISODateTimeToTime(assertedTask.Value)
						if err != nil {
							log.Printf("Task %s has invalid iso datetime value %s with the error: %s\n", assertedTask.Name, assertedTask.Value, err.Error())
						}
						// Check if current time is greater than datetime value
						if time.Now().Unix() > datetimeValue.Unix() {
							// Run task function
							log.Printf("Task %s is running\n", assertedTask.Name)
							go assertedTask.Function()
							// Update last run time
							assertedTask.LastRun = time.Now()
						}
					}
				}(task)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}