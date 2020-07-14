package utils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// IMessagingClient defines connecting, producing, and consuming messages.
type IMessagingClient interface {
	ConnectToBroker(connectionString string)
	Publish(connectionString string) error
	PubishOnQueue(msg []byte, queueName string) error
	Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error
	Close()
}

// MessagingClient encapsulates a pointer to an amqp.Connection
type MessagingClient struct {
	conn *amqp.Connection
}

// ConnectToBroker connects to a broker i.e. RabbitMQ
func (m *MessagingClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot connect to broker, connectionString not set")
	}
	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	}
}

// PublishOnQueue sends a message to a queue
func (m *MessagingClient) PublishOnQueue(body []byte, queueName string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialised. ConnectToBroker() first.")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	defer ch.Close()

	// Declare a queue that will be created if not exists
	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arugements
	)

	// Publishes a message onto the queue
	err = ch.Publish(
		"",     // Use the default exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // JSON body as []byte
		})
	log.Printf("A message was sent to queue %v: %v", queueName, body)
	return err
}
