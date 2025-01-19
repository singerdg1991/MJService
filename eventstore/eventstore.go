package eventstore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hoitek/Kit/retry"
	pb "github.com/hoitek/Maja-Service/eventstore/protobuf"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	gcred "google.golang.org/grpc/credentials/insecure"
)

// EventStore Create a new event store interface
type EventStore interface {
	PublishEvent(exchangeName string, eventNames []string, data string) (*pb.EventResponse, error)
	SubscribeEvent(eventName string) (amqp.Queue, *amqp.Channel, error)
	Close() error
}

// Create a new event store struct
type eventStore struct {
	GrpcHost           string
	GrpcPort           int
	GrpcConnection     *grpc.ClientConn
	RabbitmqHost       string
	RabbitmqPort       int
	RabbitmqUser       string
	RabbitmqPassword   string
	RabbitmqConnection *amqp.Connection
	Timout             time.Duration
}

// Default Create a default event store
var Default = EventStore(nil)

// Setup creates a new event store instance
func Setup(eventStoreGrpcHost string, eventStoreGrpcPort int, rabbitmqHost string, rabbitmqPort int, rabbitmqUser string, rabbitmqPassword string) (EventStore, error) {
	// Create a connection to the grpc server
	grpcAddress := fmt.Sprintf("%s:%d", eventStoreGrpcHost, eventStoreGrpcPort)
	grpcConn, err := retry.Get(func() (*grpc.ClientConn, error) {
		return grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(gcred.NewCredentials()))
	}, 2, 3)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Grpc of event store")

	// Create a connection to rabbitmq
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	rabbitConn, err := retry.Get(func() (*amqp.Connection, error) {
		return amqp.Dial(uri)
	}, 2, 3)
	if err != nil {
		return nil, err
	}

	// Create an event store instance
	Default = &eventStore{
		GrpcHost:           eventStoreGrpcHost,
		GrpcPort:           eventStoreGrpcPort,
		GrpcConnection:     grpcConn,
		RabbitmqHost:       rabbitmqHost,
		RabbitmqPort:       rabbitmqPort,
		RabbitmqUser:       rabbitmqUser,
		RabbitmqPassword:   rabbitmqPassword,
		RabbitmqConnection: rabbitConn,
		Timout:             5 * time.Second,
	}
	return Default, nil
}

// PublishEvent sends an event to the event store
func (e *eventStore) PublishEvent(exchangeName string, eventNames []string, data string) (*pb.EventResponse, error) {
	client := pb.NewEventServiceClient(e.GrpcConnection)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.Timout)
	defer cancel()

	// Create a random UUID for the event
	randomId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	eventId := randomId.String()

	// Send the event to the event store
	res, err := client.SendEvent(ctx, &pb.EventRequest{
		EventEntry: &pb.Event{
			EventId:      eventId,
			ExchangeName: exchangeName,
			Data:         data,
			QueueNames:   eventNames,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Return the response
	return res, nil
}

// SubscribeEvent subscribes to an event in the event store
func (e *eventStore) SubscribeEvent(eventName string) (amqp.Queue, *amqp.Channel, error) {
	ch, err := e.RabbitmqConnection.Channel()
	if err != nil {
		return amqp.Queue{}, nil, err
	}

	queue, err := ch.QueueDeclare(
		eventName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return amqp.Queue{}, nil, err
	}

	return queue, ch, nil
}

// Close closes the connection to the event store
func (e *eventStore) Close() error {
	if e.GrpcConnection != nil {
		err := e.GrpcConnection.Close()
		if err != nil {
			return err
		}
	}
	if e.RabbitmqConnection != nil {
		err := e.RabbitmqConnection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
