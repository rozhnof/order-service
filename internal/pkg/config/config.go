package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig[T any](configPath string) (T, error) {
	var cfg T

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
