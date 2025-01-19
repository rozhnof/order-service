package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	*amqp.Connection
}

const URL = "amqp://%s:%s@%s:%s/"

func NewConnection(rabbitURL string) (Connection, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		Connection: conn,
	}, nil
}
