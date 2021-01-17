package utils

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

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

// MQConfig configuration required to create a MQ connection
type MQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// NewMQ initialising the AMQP broker
func NewMQ(conf MQConfig) (*MessagingClient, error) {
	mqString := fmt.Sprintf("amqp://%s:%s@%s:%s", conf.Username, conf.Password, conf.Host, conf.Port)
	mq := &MessagingClient{}
	err := mq.ConnectToBroker(mqString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mq: %w", err)
	}
	return mq, nil
}

// Info returns AMQP info
func (m *MessagingClient) Info() string {
	return fmt.Sprintf("%s (%s)", m.conn.Properties["product"], m.conn.Properties["version"])
}

// ConnectToBroker connects to a broker i.e. RabbitMQ
func (m *MessagingClient) ConnectToBroker(connectionString string) error {
	if connectionString == "" {
		return fmt.Errorf("Cannot connect to broker, connectionString not set")
	}
	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		return fmt.Errorf("Failed to connect to AMQP compatible broker at: %s", connectionString)
	}
	return nil
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
		return fmt.Errorf("Failed to register exchange: %w", err)
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
		return fmt.Errorf("Failed to declare queue: %w", err)
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
		return fmt.Errorf("Failed to open channel: %w", err)
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
		return fmt.Errorf("Failed to register an exchange: %w", err)
	}
	q, err := ch.QueueDeclare(
		"",    // Queue name
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to register a queue: %w", err)
	}

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
		return fmt.Errorf("Failed to register a consumer: %w", err)
	}

	go consumeLoop(msgs, handlerFunc)
	return nil
}

// SubscribeToQueue registers a function for the queue.
func (m *MessagingClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to register a queue: %w", err)
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
		return fmt.Errorf("Failed to register a consumer: %w", err)
	}

	go consumeLoop(msgs, handlerFunc)
	return nil
}

// Close closes the connection to the AMQP-broker
func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}
