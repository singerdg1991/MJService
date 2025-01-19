package ports

import "github.com/streadway/amqp"

type MessageBroker interface {
	GetConnection() *amqp.Connection
	Close() error
	CreateChannel() (*amqp.Channel, error)
	CreateExchangeIfNotExist(ch *amqp.Channel, exchangeName string) error
	ConsumeQueue(queueName string) (<-chan amqp.Delivery, *amqp.Channel, error)
}
