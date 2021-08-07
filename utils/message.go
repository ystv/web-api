package utils

import (
	"fmt"

	"github.com/streadway/amqp"
)

// MQConfig configuration required to create a MQ connection
type MQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// NewMQ initialising the AMQP broker
func NewMQ(conf MQConfig) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", conf.Username, conf.Password, conf.Host, conf.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mq: %w", err)
	}
	return conn, nil
}
