package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// TODO remove logging and panicking

// IMessagingClient defines connecting, producing, and consuming messages.
type IMessagingClient interface {
	ConnectToBroker(connectionString string)
	Publish(msg []byte, exchangeName string, exchangeType string) error
	PublishOnQueue(msg []byte, queueName string) error
	PublishOnQueueWithContext(ctx context.Context, msg []byte, queueName string) error
	Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error
	Close()
	Info() string
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

// Publish publishes a message to the named exchange.
func (m *MessagingClient) Publish(body []byte, exchangeName string, exchangeType string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialised. ConnectToBroker() first.")
	}
	ch, err := m.conn.Channel() // Get a channel form the connection
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // Durable
		false, // Delete when unused
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Panicf("Failed to register exchange: %s", err.Error())
	}

	q, err := ch.QueueDeclare(
		"",    // Queue name
		false, // Durable
		false, // Delete when unused
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		log.Panicf("Failed to declare queue: %s", err.Error())
	}
	err = ch.QueueBind(
		q.Name,       // Queue name
		exchangeName, // Routing key
		exchangeName, // Exchange
		false,        // No-wait
		nil,          // Arguments
	)

	err = ch.Publish(
		exchangeName, // Exchange
		exchangeName, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			Body: body,
		})
	log.Printf("Message sent to exchange: %v", string(body))
	return err
}

// PublishOnQueueWithContext publishes the supplied body onto the named queue, passed the contenxt.
func (m *MessagingClient) PublishOnQueueWithContext(ctx context.Context, body []byte, queueName string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialised. ConnectToBroker() first.")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		buildMessage(ctx, body),
	)
	log.Printf("A message was sent to queue %v: %v", queueName, string(body))
	return err
}

func buildMessage(ctx context.Context, body []byte) amqp.Publishing {
	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	return publishing
}

// PublishOnQueue sends a message to a queue
func (m *MessagingClient) PublishOnQueue(body []byte, queueName string) error {
	return m.PublishOnQueueWithContext(nil, body, queueName)
}

// Subscribe registers a handler function for a given exchange.
func (m *MessagingClient) Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // Exchange name
		exchangeType, // Exchange type
		true,         // Durable
		false,        // Delete when unused
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register an exchange: %v", err)
	}
	log.Printf("Declared exchange, declaring queue (%s)", "")
	q, err := ch.QueueDeclare(
		"",    // Queue name
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a queue")
	}
	log.Printf("Declared queue (%d messages, %d consumers), binding to exchange (key '%s')",
		q.Messages, q.Consumers, exchangeName)

	err = ch.QueueBind(
		q.Name,       // Queue name
		exchangeName, // Bind key
		exchangeName, // Exchange
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,       // Queue name
		consumerName, // Consumer
		true,         // Auto-ack
		false,        // Exclusive
		false,        // No-local
		false,        //No-wait
		nil,          // Arguments
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go consumeLoop(msgs, handlerFunc)
	return nil
}

// SubscribeToQueue registers a function for the queue.
func (m *MessagingClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Printf("Declaring Queue (%s)", queueName)
	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Printf("Failed to register a queue: %v", err)
	}
	msgs, err := ch.Consume(
		q.Name,       // Queue name
		consumerName, // Consumer
		true,         // Auto-ack
		false,        // Exclusive
		false,        // No-local
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
	}

	go consumeLoop(msgs, handlerFunc)
	return nil
}

// Close closes the connection to the AMQP-broker
func (m *MessagingClient) Close() {
	if m.conn != nil {
		log.Printf("Closing connection to AMQP broker")
		m.conn.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}

// Info returns AMQP info
func (m *MessagingClient) Info() string {
	return fmt.Sprintf("%s (%s)", m.conn.Properties["product"], m.conn.Properties["version"])
}

// MQ global mq
var MQ IMessagingClient

// InitMessaging initialising the AMQP broker
func InitMessaging() {
	mqUsername := os.Getenv("mq_user")
	mqPassword := os.Getenv("mq_pass")
	mqHost := os.Getenv("mq_host")
	mqPort := os.Getenv("mq_port")

	mqString := fmt.Sprintf("amqp://%s:%s@%s:%s", mqUsername, mqPassword, mqHost, mqPort)
	MQ = &MessagingClient{}
	MQ.ConnectToBroker(mqString)
	log.Printf("Connected to MQ: %s@%s - %s", mqUsername, mqHost, MQ.Info())
}
