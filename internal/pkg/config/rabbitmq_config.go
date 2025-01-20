package config

type RabbitMQ struct {
	Address  string `env:"RABBITMQ_ADDRESS"          env-required:"true"`
	Port     string `env:"RABBITMQ_PORT"             env-required:"true"`
	User     string `env:"RABBITMQ_DEFAULT_USER"     env-required:"true"`
	Password string `env:"RABBITMQ_DEFAULT_PASSWORD" env-required:"true"`
}
