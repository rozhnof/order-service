package config

import "fmt"

const postgresURL = "postgres://%s:%s@%s:%s/%s?sslmode=disable"

type Postgres struct {
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Address  string `env:"POSTGRES_ADDRESS"  env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	DB       string `env:"POSTGRES_DB"       env-required:"true"`
}

func (c Postgres) ConnectionURL() string {
	return fmt.Sprintf(
		postgresURL,
		c.User,
		c.Password,
		c.Address,
		c.Port,
		c.DB,
	)
}
