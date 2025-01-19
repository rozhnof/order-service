package config

type Postgres struct {
	Address  string `env:"POSTGRES_ADDRESS"  env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DB       string `env:"POSTGRES_DB"       env-required:"true"`
}
