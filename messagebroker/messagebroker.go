package messagebroker

import (
	"fmt"
	"github.com/hoitek/Kit/retry"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	"github.com/hoitek/Maja-Service/messagebroker/ports"
	"github.com/streadway/amqp"
	"log"
)

// DefaultRabbitmqConnection is the default connection to RabbitMQ
var DefaultRabbitmqConnection *amqp.Connection

// MessageBroker is the message broker
type MessageBroker struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Connection *amqp.Connection
}

// ConnectRabbitMQ connects to RabbitMQ
func ConnectRabbitMQ(host string, port int, username string, password string) (ports.MessageBroker, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Printf("Connecting to RabbitMQ: %s", uri)
	conn, err := retry.Get(func() (*amqp.Connection, error) {
		return amqp.Dial(uri)
	}, 2, 3)

	// Connect to RabbitMQ
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to RabbitMQ: %s", uri)

	// Set RabbitMQ connection
	DefaultRabbitmqConnection = conn

	return &MessageBroker{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		Connection: conn,
	}, nil
}

// GetConnection returns the connection to RabbitMQ
func (mb *MessageBroker) GetConnection() *amqp.Connection {
	return mb.Connection
}

// Close closes the connection to RabbitMQ
func (mb *MessageBroker) Close() error {
	if mb.Connection != nil {
		if err := mb.Connection.Close(); err != nil {
			return err
		}
	}
	return nil
}

// CreateChannel creates a channel to RabbitMQ
func (mb *MessageBroker) CreateChannel() (*amqp.Channel, error) {
	ch, err := mb.Connection.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// CreateExchangeIfNotExist creates an exchange if it does not exist
func (mb *MessageBroker) CreateExchangeIfNotExist(ch *amqp.Channel, exchangeName string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

// ConsumeQueue consumes messages from a queue.
func (mb *MessageBroker) ConsumeQueue(queueName string) (<-chan amqp.Delivery, *amqp.Channel, error) {
	// Create a channel
	ch, err := mb.CreateChannel()
	if err != nil {
		return nil, nil, err
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		queueName,                           // queue name
		queueName,                           // routing key
		constants.ARRANGEMENT_EXCHANGE_NAME, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, nil, err
	}

	// Log the start of consumption
	log.Printf("Consuming messages from queue: %s", queueName)
	return msgs, ch, nil
}
