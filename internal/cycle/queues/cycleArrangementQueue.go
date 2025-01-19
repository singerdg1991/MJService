package queues

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	mbPorts "github.com/hoitek/Maja-Service/messagebroker/ports"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var isConsumerRunning bool

// CycleArrangementQueue is a struct for the RabbitMQ queue
type CycleArrangementQueue struct {
	MessageBroker     mbPorts.MessageBroker
	QueueWish         *amqp.Queue
	QueueWishResponse *amqp.Queue
}

// CycleArrangementQueueResponse is a struct for the RabbitMQ queue
type CycleArrangementQueueResponse struct {
	Error         error
	Data          string
	CorrelationID string
}

// DefaultCycleArrangementQueue is a pointer to the CycleArrangementQueue
var DefaultCycleArrangementQueue *CycleArrangementQueue

// RegisterCycleArrangementQueue register the queue for the cycle arrangement.
func RegisterCycleArrangementQueue(mb mbPorts.MessageBroker) error {
	DefaultCycleArrangementQueue = &CycleArrangementQueue{
		MessageBroker: mb,
	}

	// Create a channel
	ch, err := mb.CreateChannel()
	if err != nil {
		return err
	}

	// Create the exchange
	err = mb.CreateExchangeIfNotExist(ch, constants.ARRANGEMENT_EXCHANGE_NAME)
	if err != nil {
		if ch != nil {
			ch.Close()
		}
		return err
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		constants.ARRANGEMENT_QUEUE_WISH_NAME, // name
		true,                                  // durable
		false,                                 // delete when unused
		false,                                 // exclusive
		false,                                 // no-wait
		nil,                                   // arguments
	)
	if err != nil {
		if ch != nil {
			ch.Close()
		}
		return err
	}
	DefaultCycleArrangementQueue.QueueWish = &q

	// Declare the response queue
	qr, err := ch.QueueDeclare(
		constants.ARRANGEMENT_QUEUE_WISH_RESPONSE_NAME, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		if ch != nil {
			ch.Close()
		}
		return err
	}
	DefaultCycleArrangementQueue.QueueWishResponse = &qr

	// Bind the queue to the exchange
	err = ch.QueueBind(
		constants.ARRANGEMENT_QUEUE_WISH_NAME, // queue name
		constants.ARRANGEMENT_QUEUE_WISH_NAME, // routing key
		constants.ARRANGEMENT_EXCHANGE_NAME,   // exchange
		false,
		nil,
	)
	if err != nil {
		if ch != nil {
			ch.Close()
		}
		return err
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		constants.ARRANGEMENT_QUEUE_WISH_RESPONSE_NAME, // queue name
		constants.ARRANGEMENT_QUEUE_WISH_RESPONSE_NAME, // routing key
		constants.ARRANGEMENT_EXCHANGE_NAME,            // exchange
		false,
		nil,
	)
	if err != nil {
		if ch != nil {
			ch.Close()
		}
		return err
	}

	// Wait for response with the same correlation ID
	if !isConsumerRunning {
		go DefaultCycleArrangementQueue.ConsumeMessagesFromRequestQueueAndRespond(constants.ARRANGEMENT_QUEUE_WISH_NAME, func(data *CycleArrangementQueueResponse) (string, error) {
			log.Println("Inside ConsumeMessagesFromRequestQueueAndRespond")
			return "response text test", nil
		})
		isConsumerRunning = true
	}

	return nil
}

// Publish publishes a message to the queue.
func (q *CycleArrangementQueue) Publish(queueName string, message string, correlationID string) error {
	ch, err := q.MessageBroker.CreateChannel()
	if err != nil {
		return err
	}

	err = ch.Publish(
		constants.ARRANGEMENT_EXCHANGE_NAME, // exchange
		queueName,                           // routing key
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // Set message delivery mode to persistent
			Headers: amqp.Table{
				"correlation_id": correlationID, // Include correlation ID as a header
			},
		})
	if err != nil {
		return err
	}
	return nil
}

func (q *CycleArrangementQueue) WaitForResponseFromQueue(queueResponseName string, correlationID string) (string, error) {
	// Create a channel to receive the response message
	responseChan := make(chan string)
	defer close(responseChan)

	// Start a goroutine to listen for response messages
	go func() {
		// Create a channel to receive messages from the queue
		messages, _, err := q.MessageBroker.ConsumeQueue(queueResponseName)
		if err != nil {
			// Handle error
			responseChan <- "" // Sending empty response to signal error
			return
		}

		// Loop through messages
		for message := range messages {
			// Check if message has the expected correlation ID
			log.Println(string(message.Body))
			log.Println(message.Headers["correlation_id"].(string))
			log.Println(correlationID)
			if corrID, ok := message.Headers["correlation_id"].(string); ok && corrID == correlationID {
				// Send the response message to the channel
				responseChan <- string(message.Body)
				return
			}
		}

		// If no matching response is received within timeout, send empty response
		responseChan <- "" // Sending empty response to signal timeout
	}()

	// Wait for response or timeout
	select {
	case response := <-responseChan:
		if response == "" {
			return "", errors.New("timeout waiting for response from queue")
		}
		return response, nil
	case <-time.After(time.Second * 10): // Adjust timeout duration as needed
		return "", errors.New("timeout waiting for response from queue")
	}
}

// ConsumeMessagesFromRequestQueueAndRespond consumes messages from the request queue.
func (q *CycleArrangementQueue) ConsumeMessagesFromRequestQueueAndRespond(requestQueueName string, callback func(*CycleArrangementQueueResponse) (string, error)) {
	for {
		// Consume messages from the request queue
		msgs, _, err := q.MessageBroker.ConsumeQueue(requestQueueName)
		if err != nil {
			log.Printf("Error consuming messages from request queue: %s", err.Error())
			return
		}

		// Process incoming messages
		for message := range msgs {
			log.Printf("sdad %#v \n", message)
			// Extract correlation ID from message headers
			correlationID, ok := message.Headers["correlation_id"].(string)
			if !ok {
				log.Println("Correlation ID not found in message headers")
				continue
			}

			// Process the message
			response, err := callback(&CycleArrangementQueueResponse{
				Error:         nil,
				Data:          string(message.Body),
				CorrelationID: correlationID,
			})
			if err != nil {
				log.Printf("Error processing message: %s", err.Error())
				continue
			}

			err = q.Publish(constants.ARRANGEMENT_QUEUE_WISH_RESPONSE_NAME, response, correlationID)
			if err != nil {
				log.Printf("Error publishing response message: %s", err.Error())
				continue
			}
			message.Ack(false)
		}
	}
	log.Println("--------- After if")
}
