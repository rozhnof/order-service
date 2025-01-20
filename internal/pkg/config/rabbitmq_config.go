package config

import "fmt"

const rabbitURL = "amqp://%s:%s@%s:%s/"

type RabbitMQ struct {
	User     string `env:"RABBITMQ_DEFAULT_USER"     env-required:"true"`
	Password string `env:"RABBITMQ_DEFAULT_PASSWORD" env-required:"true"`
	Address  string `env:"RABBITMQ_ADDRESS"          env-required:"true"`
	Port     string `env:"RABBITMQ_PORT"             env-required:"true"`
}

func (c RabbitMQ) ConnectionURL() string {
	return fmt.Sprintf(
		rabbitURL,
		c.User,
		c.Password,
		c.Address,
		c.Port,
	)
}
