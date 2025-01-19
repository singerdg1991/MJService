package cycle

import (
	"github.com/hoitek/Maja-Service/internal/cycle/queues"
	mbPorts "github.com/hoitek/Maja-Service/messagebroker/ports"
)

// RegisterQueues registers the queue for the cycle arrangement.
func RegisterQueues(mb mbPorts.MessageBroker) error {
	err := queues.RegisterCycleArrangementQueue(mb)
	if err != nil {
		return err
	}
	return nil
}
